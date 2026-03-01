package user_test

// Los tests del servicio verifican la lógica de negocio de cada método.
// Se usa un mock manual del repositorio para aislar completamente el servicio
// de la base de datos.
//
// Estrategia:
//   - mockRepo implementa user.Repository con campos de tipo func.
//   - Cada subtest redefine solo los campos que necesita.
//   - Se sigue el patrón GIVEN / WHEN / THEN en cada caso.

import (
	"errors"
	"image-processing-service/internal/modules/user"
	"image-processing-service/internal/shared/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// ─────────────────────────────────────────────────────────────────────────────
// Mock del repositorio
// ─────────────────────────────────────────────────────────────────────────────

type mockRepo struct {
	CreateFn         func(u *user.User) error
	GetByEmailFn     func(email string) (*user.User, error)
	GetByIDFn        func(id string) (*user.User, error)
	GetAllFn         func() ([]*user.User, error)
	UpdateFn         func(u *user.User) error
	UpdatePasswordFn func(u *user.User) error
	DeleteFn         func(id string) error
}

func (m *mockRepo) Create(u *user.User) error                   { return m.CreateFn(u) }
func (m *mockRepo) GetByEmail(email string) (*user.User, error) { return m.GetByEmailFn(email) }
func (m *mockRepo) GetByID(id string) (*user.User, error)       { return m.GetByIDFn(id) }
func (m *mockRepo) GetAll() ([]*user.User, error)               { return m.GetAllFn() }
func (m *mockRepo) Update(u *user.User) error                   { return m.UpdateFn(u) }
func (m *mockRepo) UpdatePassword(u *user.User) error           { return m.UpdatePasswordFn(u) }
func (m *mockRepo) Delete(id string) error                      { return m.DeleteFn(id) }

// ─────────────────────────────────────────────────────────────────────────────
// GetByID
// ─────────────────────────────────────────────────────────────────────────────

func TestService_GetByID(t *testing.T) {
	repo := &mockRepo{}
	service := user.NewService(repo)
	userId := "ej55egzg4zdrs2zs6e6cxxzk"

	// ----------------------------------------------------------------
	// Caso 1: usuario no existe → ErrNotFound
	// ----------------------------------------------------------------
	t.Run("Debe retornar ErrNotFound cuando el usuario no existe", func(t *testing.T) {
		// GIVEN: el repositorio responde con ErrRecordNotFound
		repo.GetByIDFn = func(id string) (*user.User, error) {
			return nil, gorm.ErrRecordNotFound
		}

		// WHEN
		res, err := service.GetByID(userId)

		// THEN: el error de GORM se transforma en el error de dominio correcto
		assert.ErrorIs(t, err, user.ErrNotFound)
		assert.Nil(t, res)
	})

	// ----------------------------------------------------------------
	// Caso 2: error genérico del repositorio → se propaga tal cual
	// ----------------------------------------------------------------
	t.Run("Debe propagar el error genérico cuando el repositorio falla", func(t *testing.T) {
		// GIVEN: error de conexión que no es ErrRecordNotFound
		dbErr := errors.New("connection refused")
		repo.GetByIDFn = func(id string) (*user.User, error) {
			return nil, dbErr
		}

		// WHEN
		res, err := service.GetByID(userId)

		// THEN: el error llega sin transformar y NO es ErrNotFound
		assert.ErrorIs(t, err, dbErr)
		assert.NotErrorIs(t, err, user.ErrNotFound)
		assert.Nil(t, res)
	})

	// ----------------------------------------------------------------
	// Caso 3: usuario encontrado exitosamente
	// ----------------------------------------------------------------
	t.Run("Debe retornar el usuario cuando existe", func(t *testing.T) {
		// GIVEN: el repositorio devuelve un usuario válido
		expected := &user.User{ID: userId, Name: "Ana", Email: "ana@test.com"}
		repo.GetByIDFn = func(id string) (*user.User, error) {
			return expected, nil
		}

		// WHEN
		res, err := service.GetByID(userId)

		// THEN
		assert.NoError(t, err)
		assert.Equal(t, expected, res)
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// GetAll
// ─────────────────────────────────────────────────────────────────────────────

func TestService_GetAll(t *testing.T) {
	repo := &mockRepo{}
	service := user.NewService(repo)

	// ----------------------------------------------------------------
	// Caso 1: error del repositorio → se propaga
	// ----------------------------------------------------------------
	t.Run("Debe propagar el error cuando el repositorio falla", func(t *testing.T) {
		// GIVEN
		repo.GetAllFn = func() ([]*user.User, error) {
			return nil, errors.New("timeout")
		}

		// WHEN
		res, err := service.GetAll()

		// THEN
		assert.EqualError(t, err, "timeout")
		assert.Nil(t, res)
	})

	// ----------------------------------------------------------------
	// Caso 2: tabla vacía → slice vacío sin error
	// ----------------------------------------------------------------
	t.Run("Debe retornar slice vacío cuando no hay usuarios", func(t *testing.T) {
		// GIVEN
		repo.GetAllFn = func() ([]*user.User, error) {
			return []*user.User{}, nil
		}

		// WHEN
		res, err := service.GetAll()

		// THEN
		assert.NoError(t, err)
		assert.Empty(t, res)
	})

	// ----------------------------------------------------------------
	// Caso 3: múltiples usuarios → todos devueltos
	// ----------------------------------------------------------------
	t.Run("Debe retornar todos los usuarios del repositorio", func(t *testing.T) {
		// GIVEN
		usuarios := []*user.User{
			{ID: "id-1", Name: "Ana"},
			{ID: "id-2", Name: "Luis"},
		}
		repo.GetAllFn = func() ([]*user.User, error) {
			return usuarios, nil
		}

		// WHEN
		res, err := service.GetAll()

		// THEN
		assert.NoError(t, err)
		assert.Len(t, res, 2)
		assert.Equal(t, usuarios, res)
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// Update
// ─────────────────────────────────────────────────────────────────────────────

func TestService_Update(t *testing.T) {
	repo := &mockRepo{}
	service := user.NewService(repo)
	userId := "ej55egzg4zdrs2zs6e6cxxzk"

	// ----------------------------------------------------------------
	// Caso 1: usuario no existe → ErrNotFound
	// ----------------------------------------------------------------
	t.Run("Debe retornar ErrNotFound cuando el usuario no existe", func(t *testing.T) {
		// GIVEN: el repositorio responde con ErrRecordNotFound
		repo.GetByIDFn = func(id string) (*user.User, error) {
			return nil, gorm.ErrRecordNotFound
		}

		// WHEN: llamamos al servicio con cualquier payload
		req := user.UpdateUserRequest{Name: utils.Pointer("Nuevo nombre")}
		res, err := service.Update(userId, req)

		// THEN: verificamos que se transforme en el error de servicio correcto
		assert.Error(t, err)
		assert.ErrorIs(t, err, user.ErrNotFound)
		assert.Nil(t, res)
	})

	// ----------------------------------------------------------------
	// Caso 2: GetByID devuelve error genérico → se propaga sin transformar
	// ----------------------------------------------------------------
	t.Run("Debe propagar el error genérico de GetByID sin transformarlo en ErrNotFound", func(t *testing.T) {
		// GIVEN: error de conexión que no es ErrRecordNotFound
		repo.GetByIDFn = func(id string) (*user.User, error) {
			return nil, errors.New("db timeout")
		}

		// WHEN
		res, err := service.Update(userId, user.UpdateUserRequest{})

		// THEN: el error NO debe ser ErrNotFound, sino el original
		assert.EqualError(t, err, "db timeout")
		assert.NotErrorIs(t, err, user.ErrNotFound)
		assert.Nil(t, res)
	})

	// ----------------------------------------------------------------
	// Caso 3: email nuevo ya pertenece a otro usuario → ErrAlreadyExists
	// ----------------------------------------------------------------
	t.Run("Debe fallar si el email ya está en uso por otro usuario", func(t *testing.T) {
		// GIVEN: el usuario existe en la DB
		nuevoEmail := "clon@test.com"

		repo.GetByIDFn = func(id string) (*user.User, error) {
			return &user.User{ID: userId, Email: "original@test.com"}, nil
		}

		repo.GetByEmailFn = func(email string) (*user.User, error) {
			return &user.User{ID: "otro-usuario-99", Email: nuevoEmail}, nil
		}

		// WHEN: intentamos aplicar el cambio de correo
		req := user.UpdateUserRequest{Email: &nuevoEmail}
		res, err := service.Update(userId, req)

		// THEN: el servicio detecta el conflicto y retorna el error adecuado
		assert.ErrorIs(t, err, user.ErrAlreadyExists)
		assert.Nil(t, res)
	})

	// ----------------------------------------------------------------
	// Caso 4: el email nuevo es igual al actual → no consulta GetByEmail
	// ----------------------------------------------------------------
	t.Run("No debe verificar duplicado si el email enviado es igual al actual", func(t *testing.T) {
		// GIVEN: el usuario tiene el mismo email que se quiere "actualizar"
		mismoEmail := "mismo@test.com"
		repo.GetByIDFn = func(id string) (*user.User, error) {
			return &user.User{ID: userId, Email: mismoEmail}, nil
		}
		// Si GetByEmail se llamara, el test fallaría con panic
		repo.GetByEmailFn = func(email string) (*user.User, error) {
			t.Fatal("GetByEmail no debería haberse llamado cuando el email no cambia")
			return nil, nil
		}
		repo.UpdateFn = func(u *user.User) error { return nil }

		// WHEN
		req := user.UpdateUserRequest{Email: &mismoEmail}
		res, err := service.Update(userId, req)

		// THEN: no hay error y el email se mantiene igual
		assert.NoError(t, err)
		assert.Equal(t, mismoEmail, res.Email)
	})

	// ----------------------------------------------------------------
	// Caso 5: solo se actualiza el nombre → no consulta GetByEmail
	// ----------------------------------------------------------------
	t.Run("Debe actualizar solo el nombre sin consultar GetByEmail", func(t *testing.T) {
		// GIVEN: request sin campo email
		nuevoNombre := "Nombre Actualizado"
		original := &user.User{ID: userId, Name: "Viejo", Email: "fijo@test.com"}
		repo.GetByIDFn = func(id string) (*user.User, error) {
			return original, nil
		}
		repo.GetByEmailFn = func(email string) (*user.User, error) {
			t.Fatal("GetByEmail no debería haberse llamado cuando no se cambia el email")
			return nil, nil
		}
		repo.UpdateFn = func(u *user.User) error { return nil }

		// WHEN
		req := user.UpdateUserRequest{Name: utils.Pointer(nuevoNombre)}
		res, err := service.Update(userId, req)

		// THEN: nombre actualizado, email intacto
		assert.NoError(t, err)
		assert.Equal(t, nuevoNombre, res.Name)
		assert.Equal(t, "fijo@test.com", res.Email)
	})

	// ----------------------------------------------------------------
	// Caso 6: repo.Update falla → error propagado
	// ----------------------------------------------------------------
	t.Run("Debe propagar el error del repositorio cuando Update falla", func(t *testing.T) {
		// GIVEN: el usuario existe y la petición no cambia ningún campo
		usuarioOriginal := &user.User{ID: userId, Name: "X", Email: "x@x.com"}
		repo.GetByIDFn = func(id string) (*user.User, error) {
			return usuarioOriginal, nil
		}
		// WHEN: el repositorio devuelve un error en la actualización
		repo.UpdateFn = func(u *user.User) error {
			return errors.New("oops")
		}

		res, err := service.Update(userId, user.UpdateUserRequest{})

		// THEN: el error se propaga y no hay resultado
		assert.EqualError(t, err, "oops")
		assert.Nil(t, res)
	})

	// ----------------------------------------------------------------
	// Caso 7: request vacío → usuario sin cambios
	// ----------------------------------------------------------------
	t.Run("Debe devolver el usuario sin cambios si el request está vacío", func(t *testing.T) {
		// GIVEN: el usuario se recupera exitosamente y no solicitamos cambios
		usuarioOriginal := &user.User{ID: userId, Name: "Nombre", Email: "correo@ej.com"}
		repo.GetByIDFn = func(id string) (*user.User, error) {
			return usuarioOriginal, nil
		}
		// WHEN: la función Update es llamada con el mismo objeto
		repo.UpdateFn = func(u *user.User) error {
			assert.Equal(t, usuarioOriginal, u)
			return nil
		}

		res, err := service.Update(userId, user.UpdateUserRequest{})

		// THEN: no hay error y el usuario regresado es idéntico
		assert.NoError(t, err)
		assert.Equal(t, usuarioOriginal, res)
	})

	// ----------------------------------------------------------------
	// Caso 8: actualización válida de nombre y email
	// ----------------------------------------------------------------
	t.Run("Debe actualizar el nombre y el email correctamente cuando todo es válido", func(t *testing.T) {
		// GIVEN: el usuario original existe y el nuevo email está libre
		nuevoNombre := "Usuario de Prueba"
		nuevoEmail := "usuario@prueba.com"

		usuarioOriginal := &user.User{
			ID:    userId,
			Name:  "Usuario Antiguo",
			Email: "correo@anterior.com",
		}

		repo.GetByIDFn = func(id string) (*user.User, error) {
			return usuarioOriginal, nil
		}

		repo.GetByEmailFn = func(email string) (*user.User, error) {
			return nil, nil // email libre
		}

		repo.UpdateFn = func(u *user.User) error {
			return nil
		}

		// WHEN: la función de Update se llama con los valores que se pasan
		req := user.UpdateUserRequest{
			Name:  utils.Pointer(nuevoNombre),
			Email: utils.Pointer(nuevoEmail),
		}

		res, err := service.Update(userId, req)

		// THEN: el servicio devuelve los valores actualizados y mantiene el mismo ID
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, nuevoNombre, res.Name)
		assert.Equal(t, nuevoEmail, res.Email)
		assert.Equal(t, userId, res.ID)
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// UpdatePassword
// ─────────────────────────────────────────────────────────────────────────────

func TestService_UpdatePassword(t *testing.T) {
	repo := &mockRepo{}
	service := user.NewService(repo)
	userId := "ej55egzg4zdrs2zs6e6cxxzk"

	// ----------------------------------------------------------------
	// Caso 1: usuario no existe → ErrNotFound
	// ----------------------------------------------------------------
	t.Run("Debe retornar ErrNotFound cuando el usuario no existe", func(t *testing.T) {
		// GIVEN
		repo.GetByIDFn = func(id string) (*user.User, error) {
			return nil, gorm.ErrRecordNotFound
		}

		// WHEN
		req := user.UpdatePasswordUserRequest{
			CurrentPassword: "pass123",
			NewPassword:     "nueva456",
		}
		res, err := service.UpdatePassword(userId, req)

		// THEN
		assert.ErrorIs(t, err, user.ErrNotFound)
		assert.Nil(t, res)
	})

	// ----------------------------------------------------------------
	// Caso 2: GetByID devuelve error genérico → se propaga
	// ----------------------------------------------------------------
	t.Run("Debe propagar el error genérico cuando GetByID falla", func(t *testing.T) {
		// GIVEN
		repo.GetByIDFn = func(id string) (*user.User, error) {
			return nil, errors.New("db timeout")
		}

		// WHEN
		req := user.UpdatePasswordUserRequest{
			CurrentPassword: "pass123",
			NewPassword:     "nueva456",
		}
		res, err := service.UpdatePassword(userId, req)

		// THEN
		assert.EqualError(t, err, "db timeout")
		assert.Nil(t, res)
	})

	// ----------------------------------------------------------------
	// Caso 3: contraseña actual incorrecta → ErrInvalidPassword
	// ----------------------------------------------------------------
	t.Run("Debe retornar ErrInvalidPassword cuando la contraseña no coincide", func(t *testing.T) {
		// GIVEN: usuario con contraseña hasheada conocida
		hashedPass, _ := utils.HashPassword("correcta123")
		repo.GetByIDFn = func(id string) (*user.User, error) {
			return &user.User{ID: userId, Password: hashedPass}, nil
		}

		// WHEN: se envía la contraseña incorrecta
		req := user.UpdatePasswordUserRequest{
			CurrentPassword: "incorrecta999",
			NewPassword:     "nueva456",
		}
		res, err := service.UpdatePassword(userId, req)

		// THEN
		assert.ErrorIs(t, err, user.ErrInvalidPassword)
		assert.Nil(t, res)
	})

	// ----------------------------------------------------------------
	// Caso 4: repo.UpdatePassword falla → error propagado
	// ----------------------------------------------------------------
	t.Run("Debe propagar el error cuando UpdatePassword del repositorio falla", func(t *testing.T) {
		// GIVEN
		hashedPass, _ := utils.HashPassword("correcta123")
		repo.GetByIDFn = func(id string) (*user.User, error) {
			return &user.User{ID: userId, Password: hashedPass}, nil
		}
		repo.UpdatePasswordFn = func(u *user.User) error {
			return errors.New("db error")
		}

		// WHEN
		req := user.UpdatePasswordUserRequest{
			CurrentPassword: "correcta123",
			NewPassword:     "nueva456",
		}
		res, err := service.UpdatePassword(userId, req)

		// THEN
		assert.EqualError(t, err, "db error")
		assert.Nil(t, res)
	})

	// ----------------------------------------------------------------
	// Caso 5: flujo exitoso → contraseña actualizada y hasheada
	// ----------------------------------------------------------------
	t.Run("Debe actualizar la contraseña correctamente cuando todo es válido", func(t *testing.T) {
		// GIVEN: usuario con contraseña hasheada conocida
		hashedPass, _ := utils.HashPassword("correcta123")
		repo.GetByIDFn = func(id string) (*user.User, error) {
			return &user.User{ID: userId, Password: hashedPass}, nil
		}
		repo.UpdatePasswordFn = func(u *user.User) error {
			return nil
		}

		// WHEN
		req := user.UpdatePasswordUserRequest{
			CurrentPassword: "correcta123",
			NewPassword:     "nueva456",
		}
		res, err := service.UpdatePassword(userId, req)

		// THEN: la contraseña almacenada ya no es la original y verifica la nueva
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotEqual(t, hashedPass, res.Password)
		assert.True(t, utils.CheckPasswordHash("nueva456", res.Password))
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// Delete
// ─────────────────────────────────────────────────────────────────────────────

func TestService_Delete(t *testing.T) {
	repo := &mockRepo{}
	service := user.NewService(repo)
	userId := "ej55egzg4zdrs2zs6e6cxxzk"

	// ----------------------------------------------------------------
	// Caso 1: usuario no existe → ErrNotFound
	// ----------------------------------------------------------------
	t.Run("Debe retornar ErrNotFound cuando el usuario no existe", func(t *testing.T) {
		// GIVEN
		repo.DeleteFn = func(id string) error {
			return gorm.ErrRecordNotFound
		}

		// WHEN
		err := service.Delete(userId)

		// THEN
		assert.ErrorIs(t, err, user.ErrNotFound)
	})

	// ----------------------------------------------------------------
	// Caso 2: error genérico del repositorio → se propaga
	// ----------------------------------------------------------------
	t.Run("Debe propagar el error genérico cuando el repositorio falla", func(t *testing.T) {
		// GIVEN
		repo.DeleteFn = func(id string) error {
			return errors.New("constraint violation")
		}

		// WHEN
		err := service.Delete(userId)

		// THEN
		assert.EqualError(t, err, "constraint violation")
	})

	// ----------------------------------------------------------------
	// Caso 3: eliminación exitosa
	// ----------------------------------------------------------------
	t.Run("Debe eliminar el usuario sin error cuando existe", func(t *testing.T) {
		// GIVEN
		repo.DeleteFn = func(id string) error {
			return nil
		}

		// WHEN
		err := service.Delete(userId)

		// THEN
		assert.NoError(t, err)
	})
}
