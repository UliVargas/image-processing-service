package user_test

// Los tests de repositorio verifican que las operaciones de base de datos
// funcionan correctamente. Se usa SQLite en memoria como base de datos real,
// lo que permite ejecutar las queries sin necesitar un servidor externo.
//
// Estrategia:
//   - Cada test crea una base de datos SQLite en memoria completamente nueva.
//   - GORM migra el esquema automáticamente antes de cada test.
//   - Se insertan datos de prueba directamente con GORM (sin mocks).
//   - Se verifica el resultado de las operaciones del repositorio.
//
// Por qué SQLite en lugar de sqlmock:
//   - No hay que declarar queries esperadas: el test es más legible.
//   - Se prueban las queries reales, no una representación de ellas.
//   - Es el enfoque estándar de la comunidad Go para tests de repositorio.
//   - Cada test es completamente aislado: la DB se descarta al terminar.

import (
	"errors"
	"fmt"
	"image-processing-service/internal/modules/user"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// newMemoryDB crea una base de datos SQLite en memoria con el esquema del
// módulo user ya migrado. Se llama al inicio de cada subtest.
func newMemoryDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err, "no se pudo abrir la base de datos en memoria")

	// AutoMigrate crea la tabla users con todas sus columnas y restricciones
	err = db.AutoMigrate(&user.User{})
	require.NoError(t, err, "no se pudo migrar el esquema")

	return db
}

// ─────────────────────────────────────────────────────────────────────────────
// GetByEmail
// ─────────────────────────────────────────────────────────────────────────────

func TestRepository_GetByEmail(t *testing.T) {
	// ----------------------------------------------------------------
	// Caso 1: email no existe en la base de datos
	// ----------------------------------------------------------------
	t.Run("Debe retornar nil cuando el email no existe", func(t *testing.T) {
		// GIVEN: base de datos vacía
		db := newMemoryDB(t)
		repo := user.NewRepository(db)

		// WHEN
		result, err := repo.GetByEmail("noexiste@test.com")

		// THEN: nil, nil (comportamiento especial del repositorio para "no encontrado")
		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	// ----------------------------------------------------------------
	// Caso 2: email encontrado exitosamente
	// ----------------------------------------------------------------
	t.Run("Debe retornar el usuario cuando el email existe", func(t *testing.T) {
		// GIVEN: usuario insertado directamente en la DB
		db := newMemoryDB(t)
		repo := user.NewRepository(db)

		esperado := &user.User{ID: "id-1", Name: "Ana García", Email: "ana@test.com", Password: "hash"}
		require.NoError(t, db.Create(esperado).Error)

		// WHEN
		result, err := repo.GetByEmail("ana@test.com")

		// THEN
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "ana@test.com", result.Email)
		assert.Equal(t, "Ana García", result.Name)
	})

	// ----------------------------------------------------------------
	// Caso 3: usuario eliminado (soft delete) → aún se encuentra
	// ----------------------------------------------------------------
	t.Run("Debe retornar el usuario aunque haya sido eliminado (soft delete)", func(t *testing.T) {
		// GIVEN: usuario creado y luego eliminado con soft delete
		db := newMemoryDB(t)
		repo := user.NewRepository(db)

		u := &user.User{ID: "id-2", Name: "Luis", Email: "luis@test.com", Password: "hash"}
		require.NoError(t, db.Create(u).Error)
		require.NoError(t, db.Delete(u).Error) // soft delete

		// WHEN: GetByEmail usa Unscoped, por lo que encuentra registros eliminados
		result, err := repo.GetByEmail("luis@test.com")

		// THEN: el repositorio usa Unscoped para buscar por email
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "luis@test.com", result.Email)
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// GetByID
// ─────────────────────────────────────────────────────────────────────────────

func TestRepository_GetByID(t *testing.T) {
	// ----------------------------------------------------------------
	// Caso 1: ID no existe → gorm.ErrRecordNotFound
	// ----------------------------------------------------------------
	t.Run("Debe retornar ErrRecordNotFound cuando el ID no existe", func(t *testing.T) {
		// GIVEN: base de datos vacía
		db := newMemoryDB(t)
		repo := user.NewRepository(db)

		// WHEN
		result, err := repo.GetByID("id-inexistente")

		// THEN
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
		assert.Nil(t, result)
	})

	// ----------------------------------------------------------------
	// Caso 2: ID encontrado exitosamente
	// ----------------------------------------------------------------
	t.Run("Debe retornar el usuario cuando el ID existe", func(t *testing.T) {
		// GIVEN: usuario insertado directamente en la DB
		db := newMemoryDB(t)
		repo := user.NewRepository(db)

		esperado := &user.User{ID: "id-valido", Name: "Luis Pérez", Email: "luis@test.com", Password: "hash"}
		require.NoError(t, db.Create(esperado).Error)

		// WHEN
		result, err := repo.GetByID("id-valido")

		// THEN
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "id-valido", result.ID)
		assert.Equal(t, "Luis Pérez", result.Name)
	})

	// ----------------------------------------------------------------
	// Caso 3: usuario eliminado (soft delete) → no se encuentra
	// ----------------------------------------------------------------
	t.Run("Debe retornar ErrRecordNotFound cuando el usuario fue eliminado", func(t *testing.T) {
		// GIVEN: usuario creado y luego eliminado con soft delete
		db := newMemoryDB(t)
		repo := user.NewRepository(db)

		u := &user.User{ID: "id-eliminado", Name: "María", Email: "maria@test.com", Password: "hash"}
		require.NoError(t, db.Create(u).Error)
		require.NoError(t, db.Delete(u).Error) // soft delete

		// WHEN: GetByID no usa Unscoped, por lo que no encuentra eliminados
		result, err := repo.GetByID("id-eliminado")

		// THEN
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
		assert.Nil(t, result)
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// GetAll
// ─────────────────────────────────────────────────────────────────────────────

func TestRepository_GetAll(t *testing.T) {
	// ----------------------------------------------------------------
	// Caso 1: tabla vacía → slice vacío con total=0
	// ----------------------------------------------------------------
	t.Run("Debe retornar slice vacío cuando no hay usuarios", func(t *testing.T) {
		// GIVEN: base de datos vacía
		db := newMemoryDB(t)
		repo := user.NewRepository(db)

		// WHEN
		result, total, err := repo.GetAll(1, 10)

		// THEN
		assert.NoError(t, err)
		assert.Empty(t, result)
		assert.Equal(t, int64(0), total)
	})

	// ----------------------------------------------------------------
	// Caso 2: múltiples usuarios → todos devueltos con total correcto
	// ----------------------------------------------------------------
	t.Run("Debe retornar todos los usuarios de la base de datos", func(t *testing.T) {
		// GIVEN: dos usuarios insertados
		db := newMemoryDB(t)
		repo := user.NewRepository(db)

		require.NoError(t, db.Create(&user.User{ID: "id-1", Name: "Ana", Email: "ana@test.com", Password: "h1"}).Error)
		require.NoError(t, db.Create(&user.User{ID: "id-2", Name: "Luis", Email: "luis@test.com", Password: "h2"}).Error)

		// WHEN
		result, total, err := repo.GetAll(1, 10)

		// THEN
		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, int64(2), total)
	})

	// ----------------------------------------------------------------
	// Caso 3: usuarios eliminados (soft delete) → no se incluyen en total ni en data
	// ----------------------------------------------------------------
	t.Run("No debe incluir usuarios eliminados en el resultado", func(t *testing.T) {
		// GIVEN: un usuario activo y uno eliminado
		db := newMemoryDB(t)
		repo := user.NewRepository(db)

		activo := &user.User{ID: "id-activo", Name: "Ana", Email: "ana@test.com", Password: "h1"}
		eliminado := &user.User{ID: "id-eliminado", Name: "Luis", Email: "luis@test.com", Password: "h2"}
		require.NoError(t, db.Create(activo).Error)
		require.NoError(t, db.Create(eliminado).Error)
		require.NoError(t, db.Delete(eliminado).Error) // soft delete

		// WHEN
		result, total, err := repo.GetAll(1, 10)

		// THEN: solo el usuario activo aparece
		require.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, int64(1), total)
		assert.Equal(t, "id-activo", result[0].ID)
	})

	// ----------------------------------------------------------------
	// Caso 4: LIMIT y OFFSET funcionan correctamente
	// ----------------------------------------------------------------
	t.Run("Debe respetar el LIMIT y OFFSET según page y limit", func(t *testing.T) {
		// GIVEN: 5 usuarios insertados
		db := newMemoryDB(t)
		repo := user.NewRepository(db)

		for i := 1; i <= 5; i++ {
			require.NoError(t, db.Create(&user.User{
				ID:       fmt.Sprintf("id-%d", i),
				Name:     fmt.Sprintf("User%d", i),
				Email:    fmt.Sprintf("u%d@test.com", i),
				Password: "h",
			}).Error)
		}

		// WHEN: página 2 con limite 2 → usuarios 3 y 4
		result, total, err := repo.GetAll(2, 2)

		// THEN
		require.NoError(t, err)
		assert.Equal(t, int64(5), total) // total siempre refleja el real
		assert.Len(t, result, 2)         // solo 2 por página
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// Create
// ─────────────────────────────────────────────────────────────────────────────

func TestRepository_Create(t *testing.T) {
	// ----------------------------------------------------------------
	// Caso 1: creación exitosa
	// ----------------------------------------------------------------
	t.Run("Debe crear el usuario sin error cuando los datos son válidos", func(t *testing.T) {
		// GIVEN: base de datos vacía
		db := newMemoryDB(t)
		repo := user.NewRepository(db)

		u := &user.User{ID: "id-nuevo", Name: "María", Email: "maria@test.com", Password: "hash"}

		// WHEN
		err := repo.Create(u)

		// THEN: sin error y el usuario existe en la DB
		require.NoError(t, err)

		var guardado user.User
		require.NoError(t, db.First(&guardado, "id = ?", "id-nuevo").Error)
		assert.Equal(t, "María", guardado.Name)
		assert.Equal(t, "maria@test.com", guardado.Email)
	})

	// ----------------------------------------------------------------
	// Caso 2: email duplicado → error de constraint
	// ----------------------------------------------------------------
	t.Run("Debe retornar error cuando el email ya existe", func(t *testing.T) {
		// GIVEN: usuario con ese email ya existe
		db := newMemoryDB(t)
		repo := user.NewRepository(db)

		require.NoError(t, db.Create(&user.User{
			ID: "id-1", Name: "Ana", Email: "duplicado@test.com", Password: "hash",
		}).Error)

		// WHEN: intentamos crear otro usuario con el mismo email
		u2 := &user.User{ID: "id-2", Name: "Otro", Email: "duplicado@test.com", Password: "hash"}
		err := repo.Create(u2)

		// THEN: la restricción de unicidad del email genera un error
		assert.Error(t, err)
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// Update
// ─────────────────────────────────────────────────────────────────────────────

func TestRepository_Update(t *testing.T) {
	// ----------------------------------------------------------------
	// Caso 1: actualización exitosa
	// ----------------------------------------------------------------
	t.Run("Debe actualizar el usuario sin error cuando los datos son válidos", func(t *testing.T) {
		// GIVEN: usuario existente en la DB
		db := newMemoryDB(t)
		repo := user.NewRepository(db)

		original := &user.User{ID: "id-1", Name: "Original", Email: "original@test.com", Password: "hash"}
		require.NoError(t, db.Create(original).Error)

		// WHEN: actualizamos nombre y email
		original.Name = "Actualizado"
		original.Email = "actualizado@test.com"
		err := repo.Update(original)

		// THEN: sin error y los cambios persisten en la DB
		require.NoError(t, err)

		var guardado user.User
		require.NoError(t, db.First(&guardado, "id = ?", "id-1").Error)
		assert.Equal(t, "Actualizado", guardado.Name)
		assert.Equal(t, "actualizado@test.com", guardado.Email)
	})

	// ----------------------------------------------------------------
	// Caso 2: Update no modifica la contraseña
	// ----------------------------------------------------------------
	t.Run("No debe modificar la contraseña al actualizar nombre y email", func(t *testing.T) {
		// GIVEN: usuario con contraseña conocida
		db := newMemoryDB(t)
		repo := user.NewRepository(db)

		original := &user.User{ID: "id-1", Name: "Ana", Email: "ana@test.com", Password: "hash-original"}
		require.NoError(t, db.Create(original).Error)

		// WHEN: actualizamos nombre pero intentamos cambiar la contraseña
		original.Name = "Ana Actualizada"
		original.Password = "nueva-contraseña"
		err := repo.Update(original)

		// THEN: el repositorio omite el campo password en Update
		require.NoError(t, err)

		var guardado user.User
		require.NoError(t, db.First(&guardado, "id = ?", "id-1").Error)
		assert.Equal(t, "Ana Actualizada", guardado.Name)
		assert.Equal(t, "hash-original", guardado.Password) // contraseña sin cambios
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// UpdatePassword
// ─────────────────────────────────────────────────────────────────────────────

func TestRepository_UpdatePassword(t *testing.T) {
	// ----------------------------------------------------------------
	// Caso 1: actualización de contraseña exitosa
	// ----------------------------------------------------------------
	t.Run("Debe actualizar solo la contraseña sin modificar otros campos", func(t *testing.T) {
		// GIVEN: usuario existente
		db := newMemoryDB(t)
		repo := user.NewRepository(db)

		original := &user.User{ID: "id-1", Name: "Ana", Email: "ana@test.com", Password: "hash-viejo"}
		require.NoError(t, db.Create(original).Error)

		// WHEN: actualizamos la contraseña
		original.Password = "hash-nuevo"
		err := repo.UpdatePassword(original)

		// THEN: solo la contraseña cambió
		require.NoError(t, err)

		var guardado user.User
		require.NoError(t, db.First(&guardado, "id = ?", "id-1").Error)
		assert.Equal(t, "hash-nuevo", guardado.Password)
		assert.Equal(t, "Ana", guardado.Name)           // nombre sin cambios
		assert.Equal(t, "ana@test.com", guardado.Email) // email sin cambios
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// Delete
// ─────────────────────────────────────────────────────────────────────────────

func TestRepository_Delete(t *testing.T) {
	// ----------------------------------------------------------------
	// Caso 1: ID no existe → ErrRecordNotFound
	// ----------------------------------------------------------------
	t.Run("Debe retornar ErrRecordNotFound cuando el ID no existe", func(t *testing.T) {
		// GIVEN: base de datos vacía
		db := newMemoryDB(t)
		repo := user.NewRepository(db)

		// WHEN
		err := repo.Delete("id-inexistente")

		// THEN
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})

	// ----------------------------------------------------------------
	// Caso 2: eliminación exitosa (soft delete)
	// ----------------------------------------------------------------
	t.Run("Debe eliminar el usuario sin error cuando existe", func(t *testing.T) {
		// GIVEN: usuario existente
		db := newMemoryDB(t)
		repo := user.NewRepository(db)

		u := &user.User{ID: "id-existente", Name: "Ana", Email: "ana@test.com", Password: "hash"}
		require.NoError(t, db.Create(u).Error)

		// WHEN
		err := repo.Delete("id-existente")

		// THEN: sin error y el usuario ya no se encuentra con GetByID
		require.NoError(t, err)

		var encontrado user.User
		result := db.First(&encontrado, "id = ?", "id-existente")
		assert.ErrorIs(t, result.Error, gorm.ErrRecordNotFound)
	})

	// ----------------------------------------------------------------
	// Caso 3: doble eliminación → ErrRecordNotFound en el segundo intento
	// ----------------------------------------------------------------
	t.Run("Debe retornar ErrRecordNotFound al intentar eliminar un usuario ya eliminado", func(t *testing.T) {
		// GIVEN: usuario eliminado previamente
		db := newMemoryDB(t)
		repo := user.NewRepository(db)

		u := &user.User{ID: "id-ya-eliminado", Name: "Luis", Email: "luis@test.com", Password: "hash"}
		require.NoError(t, db.Create(u).Error)
		require.NoError(t, repo.Delete("id-ya-eliminado")) // primera eliminación

		// WHEN: intentamos eliminar de nuevo
		err := repo.Delete("id-ya-eliminado")

		// THEN: el repositorio detecta que no hay filas afectadas
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})
}

// Aseguramos que errors se importa
var _ = errors.New
