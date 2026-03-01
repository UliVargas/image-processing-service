# Guía de Testing

Esta guía es el documento canónico que define la estrategia de pruebas del
proyecto. Está escrita de forma prescriptiva: cada regla incluye su
justificación. Un desarrollador sin contexto previo debe poder escribir un
archivo de test indistinguible en estilo y calidad de los existentes
siguiendo únicamente este documento.

---

## Índice

1. [Filosofía y objetivos](#1-filosofía-y-objetivos)
2. [Estructura y organización](#2-estructura-y-organización)
3. [Convenciones de nomenclatura](#3-convenciones-de-nomenclatura)
4. [El patrón GIVEN / WHEN / THEN](#4-el-patrón-given--when--then)
5. [Mocks manuales: el patrón func-field](#5-mocks-manuales-el-patrón-func-field)
6. [Aislamiento de dependencias](#6-aislamiento-de-dependencias)
7. [Subtests individuales vs tests con tabla](#7-subtests-individuales-vs-tests-con-tabla)
8. [Pruebas de casos de error](#8-pruebas-de-casos-de-error)
9. [Aserciones sin librerías de terceros](#9-aserciones-sin-librerías-de-terceros)
10. [Setup y teardown](#10-setup-y-teardown)
11. [Tests de integración vs tests unitarios](#11-tests-de-integración-vs-tests-unitarios)
12. [Qué siempre, qué a veces, qué nunca probar](#12-qué-siempre-qué-a-veces-qué-nunca-probar)
13. [Anti-patrones a evitar](#13-anti-patrones-a-evitar)
14. [Ejemplo completo de archivo de test unitario](#14-ejemplo-completo-de-archivo-de-test-unitario)

---

## 1. Filosofía y objetivos

### Principio central

Los tests no verifican que el código *existe*; verifican que el código
*se comporta correctamente* ante cada escenario posible. Un test que
siempre pasa no aporta valor.

### Objetivos en orden de prioridad

1. **Detectar regresiones.** Un cambio que rompe comportamiento existente
   debe fallar en CI antes de llegar a producción.
2. **Documentar el comportamiento.** Los tests son la especificación
   ejecutable del sistema. Leer los nombres de los subtests debe ser
   suficiente para entender qué hace cada función.
3. **Facilitar el refactoring.** Los tests deben acoplarse al
   *comportamiento observable*, no a los detalles de implementación.
   Si se puede refactorizar el código sin tocar los tests, los tests
   están bien escritos.
4. **Dar confianza rápida.** La suite completa debe ejecutarse en
   segundos. Los tests lentos se ignoran.

### Qué no son los tests

- No son una métrica de cobertura. El 100 % de cobertura con aserciones
  vacías es peor que el 60 % con aserciones precisas.
- No son documentación de la implementación interna. Si un test falla
  porque se renombró una variable privada, el test está mal escrito.

---

## 2. Estructura y organización

### Ubicación de los archivos

Cada archivo de test vive en el mismo directorio que el código que prueba
y lleva el sufijo `_test.go`:

```
internal/modules/example/
├── service.go
├── service_test.go      ← tests unitarios del servicio
├── repository.go
├── repository_test.go   ← tests de integración del repositorio
├── handler.go
└── handler_test.go      ← tests HTTP del handler
```

**Justificación:** Go compila los archivos `_test.go` solo durante
`go test`, por lo que no aumentan el tamaño del binario de producción.
Mantenerlos junto al código facilita la navegación y hace evidente
cuándo falta cobertura.

### Declaración del paquete

Todos los archivos de test usan el paquete con sufijo `_test`:

```go
// ✅ Correcto
package example_test

// ❌ Incorrecto
package example
```

**Justificación:** El paquete `example_test` es una caja negra: solo
puede acceder a los identificadores exportados del paquete `example`.
Esto garantiza que los tests validan la API pública, no los detalles
internos. Si algo no se puede probar desde fuera, es una señal de que
la API necesita revisión, no de que el test deba usar el paquete interno.

### Un `TestXxx` por método público

Cada función o método público tiene su propia función `TestXxx`. Los
escenarios de ese método se organizan como subtests dentro de ella.

```go
func TestService_GetByID(t *testing.T) { ... }
func TestService_Update(t *testing.T)  { ... }
func TestService_Delete(t *testing.T)  { ... }
```

**Justificación:** Agrupa los casos relacionados, permite ejecutar solo
los tests de un método con `-run TestService_GetByID`, y hace que el
reporte de fallos sea inmediatamente localizable.

---

## 3. Convenciones de nomenclatura

### Funciones de test

El formato es `Test{Capa}_{Método}`:

| Capa | Prefijo |
|------|---------|
| Servicio | `TestService_` |
| Repositorio | `TestRepository_` |
| Handler HTTP | `TestHandler_` |

Ejemplos:

```go
func TestService_GetByID(t *testing.T)        {}
func TestRepository_Create(t *testing.T)      {}
func TestHandler_UpdatePassword(t *testing.T) {}
```

### Subtests

Los nombres de subtests son frases en español que describen el escenario
completo. Deben poder leerse como una especificación:

```go
// ✅ Correcto: describe el escenario y el resultado esperado
t.Run("Debe retornar ErrNotFound cuando el usuario no existe", ...)
t.Run("Debe propagar el error del repositorio cuando Update falla", ...)
t.Run("Debe actualizar el nombre y el email cuando todo es válido", ...)

// ❌ Incorrecto: vago, no describe el resultado
t.Run("test error", ...)
t.Run("caso 1", ...)
t.Run("success", ...)
```

**Justificación:** El nombre del subtest aparece en el reporte de fallos.
Un nombre descriptivo elimina la necesidad de leer el código del test
para entender qué falló.

### Variables dentro de los tests

| Variable | Convención |
|----------|------------|
| Instancia del mock | `repo`, `svc` |
| Instancia bajo prueba | `service`, `handler` |
| ID de prueba | `userId`, `validID` |
| Request de entrada | `req` |
| Resultado | `res`, `result` |
| Error | `err` |
| Datos de entrada | nombre descriptivo (`nuevoEmail`, `hashedPass`) |

---

## 4. El patrón GIVEN / WHEN / THEN

Cada subtest se divide en tres secciones marcadas con comentarios:

```go
t.Run("Debe retornar ErrNotFound cuando el usuario no existe", func(t *testing.T) {
    // GIVEN: estado inicial del sistema (configurar mocks y datos de entrada)
    repo.GetByIDFn = func(id string) (*SomeType, error) {
        return nil, ErrRecordNotFound
    }

    // WHEN: acción que se prueba (una sola llamada al código bajo prueba)
    res, err := service.GetByID("algún-id")

    // THEN: verificaciones del resultado esperado
    assert.ErrorIs(t, err, ErrNotFound)
    assert.Nil(t, res)
})
```

### Reglas de cada sección

**GIVEN**
- Configura únicamente lo que el test necesita. No configures mocks que
  no se van a invocar en este escenario.
- Si el test necesita datos de entrada complejos, créalos aquí con
  nombres descriptivos.

**WHEN**
- Contiene exactamente **una** llamada al código bajo prueba.
- No hay lógica condicional en esta sección.

**THEN**
- Verifica el resultado y el error.
- Verifica efectos secundarios observables (por ejemplo, que un mock
  fue llamado con los argumentos correctos).
- No hay lógica condicional en esta sección.

### Por qué una sola llamada en WHEN

Si el WHEN tiene dos llamadas, el test está probando dos cosas a la vez.
Cuando falla, no sabes cuál de las dos causó el fallo. Divide el test.

---

## 5. Mocks manuales: el patrón func-field

### Definición

Un mock manual es una implementación de una interfaz donde cada método
delega en un campo de tipo `func`. Esto permite redefinir el
comportamiento de cada método en cada subtest sin crear una nueva
instancia del mock.

```go
// Definición del mock (una sola vez por archivo de test)
type mockRepo struct {
    GetByIDFn func(id string) (*SomeType, error)
    UpdateFn  func(u *SomeType) error
    DeleteFn  func(id string) error
}

// Implementación de la interfaz (delega en el campo Fn)
func (m *mockRepo) GetByID(id string) (*SomeType, error) { return m.GetByIDFn(id) }
func (m *mockRepo) Update(u *SomeType) error             { return m.UpdateFn(u) }
func (m *mockRepo) Delete(id string) error               { return m.DeleteFn(id) }
```

### Uso en los tests

```go
func TestService_Update(t *testing.T) {
    repo := &mockRepo{}          // una sola instancia para todos los subtests
    service := NewService(repo)  // inyección de dependencia

    t.Run("Caso A", func(t *testing.T) {
        // GIVEN: redefinimos solo los métodos que este caso necesita
        repo.GetByIDFn = func(id string) (*SomeType, error) {
            return nil, ErrNotFound
        }
        // ...
    })

    t.Run("Caso B", func(t *testing.T) {
        // GIVEN: redefinimos para este caso (sobreescribe el anterior)
        repo.GetByIDFn = func(id string) (*SomeType, error) {
            return &SomeType{ID: id}, nil
        }
        repo.UpdateFn = func(u *SomeType) error {
            return nil
        }
        // ...
    })
}
```

### Por qué mocks manuales y no una librería

- **Sin dependencias externas.** Las librerías de mocking añaden
  complejidad al módulo y tienen su propia curva de aprendizaje.
- **Legibilidad.** El comportamiento del mock es código Go normal,
  visible directamente en el test.
- **Control total.** Puedes hacer aserciones dentro del mock para
  verificar que se llamó con los argumentos correctos.
- **Compilación estricta.** Si la interfaz cambia, el compilador señala
  exactamente qué mocks necesitan actualizarse.

### Aserciones dentro del mock

Cuando necesitas verificar que el código bajo prueba llama a un método
con los argumentos correctos, haz la aserción dentro del campo `Fn`:

```go
repo.UpdateFn = func(u *SomeType) error {
    // Verificamos que el objeto llegó con los valores esperados
    assert.Equal(t, "nuevo-nombre", u.Name)
    assert.Equal(t, "nuevo@email.com", u.Email)
    return nil
}
```

### Detectar llamadas que no deberían ocurrir

Si un método no debería ser invocado en un escenario, usa `t.Fatal`
dentro del campo `Fn` para que el test falle inmediatamente si se llama:

```go
repo.GetByEmailFn = func(email string) (*SomeType, error) {
    t.Fatal("GetByEmail no debería haberse llamado en este escenario")
    return nil, nil
}
```

---

## 6. Aislamiento de dependencias

### Principio

El código bajo prueba no debe tener acceso a ningún recurso externo
(base de datos, red, sistema de archivos, reloj del sistema) durante
los tests unitarios.

### Cómo lograrlo

**Inyección de dependencias:** todas las dependencias externas se pasan
como interfaces al constructor. Nunca se instancian dentro de la función
que se prueba.

```go
// ✅ Correcto: la dependencia se inyecta
func NewService(repo Repository) Service {
    return &service{repo: repo}
}

// ❌ Incorrecto: la dependencia se crea internamente
func NewService() Service {
    return &service{repo: NewPostgresRepository(db)}
}
```

**Interfaces pequeñas:** define interfaces con solo los métodos que el
consumidor necesita. Esto hace que los mocks sean más simples y que el
código sea más fácil de probar.

**Sin variables globales mutables:** las variables globales que cambian
de estado hacen que los tests sean no deterministas. Si necesitas
compartir estado entre tests, usa `t.Helper()` y funciones de setup
explícitas.

---

## 7. Subtests individuales vs tests con tabla

### Cuándo usar subtests individuales (`t.Run`)

Usa subtests individuales cuando los escenarios tienen **configuraciones
de mock distintas** o cuando el GIVEN de cada caso es significativamente
diferente:

```go
func TestService_Update(t *testing.T) {
    repo := &mockRepo{}
    service := NewService(repo)

    t.Run("Debe retornar ErrNotFound cuando el usuario no existe", func(t *testing.T) {
        repo.GetByIDFn = func(id string) (*SomeType, error) {
            return nil, ErrRecordNotFound
        }
        _, err := service.Update("id", UpdateRequest{})
        assert.ErrorIs(t, err, ErrNotFound)
    })

    t.Run("Debe propagar el error genérico del repositorio", func(t *testing.T) {
        repo.GetByIDFn = func(id string) (*SomeType, error) {
            return nil, errors.New("timeout")
        }
        _, err := service.Update("id", UpdateRequest{})
        assert.EqualError(t, err, "timeout")
    })
}
```

### Cuándo usar tests con tabla

Usa tests con tabla cuando los escenarios tienen la **misma estructura**
(mismo GIVEN, mismo WHEN) y solo varían los datos de entrada y el
resultado esperado:

```go
func TestValidateEmail(t *testing.T) {
    casos := []struct {
        nombre   string
        email    string
        esperado bool
    }{
        {"email válido", "usuario@dominio.com", true},
        {"sin arroba", "usuariodominio.com", false},
        {"sin dominio", "usuario@", false},
        {"vacío", "", false},
    }

    for _, tc := range casos {
        t.Run(tc.nombre, func(t *testing.T) {
            resultado := ValidateEmail(tc.email)
            assert.Equal(t, tc.esperado, resultado)
        })
    }
}
```

### Regla de decisión

Si necesitas cambiar la configuración del mock entre casos, usa subtests
individuales. Si solo cambian los datos de entrada y el resultado
esperado, usa tabla.

---

## 8. Pruebas de casos de error

### Cobertura mínima de errores

Para cada función que puede devolver un error, se deben probar:

1. **Cada tipo de error distinto** que la función puede devolver.
2. **La transformación de errores:** si la función convierte un error de
   infraestructura en un error de dominio, verifica ambos lados.
3. **La propagación de errores:** si la función no transforma el error,
   verifica que llega sin modificar.

### Cómo construir errores esperados

```go
// Para errores de dominio definidos como variables del paquete:
assert.ErrorIs(t, err, ErrNotFound)
assert.ErrorIs(t, err, ErrAlreadyExists)

// Para errores con mensaje específico:
assert.EqualError(t, err, "connection refused")

// Para verificar que el error NO es de un tipo específico:
assert.NotErrorIs(t, err, ErrNotFound)
```

### Cómo evitar falsos positivos

Un falso positivo es un test que pasa aunque el código esté roto.
Las causas más comunes son:

**Aserción de error sin verificar el tipo:**
```go
// ❌ Pasa aunque err sea cualquier error
assert.Error(t, err)

// ✅ Verifica que es exactamente el error esperado
assert.ErrorIs(t, err, ErrNotFound)
```

**No verificar que el resultado es nil cuando hay error:**
```go
// ❌ Incompleto: no verifica que res sea nil
assert.ErrorIs(t, err, ErrNotFound)

// ✅ Completo: verifica error Y resultado
assert.ErrorIs(t, err, ErrNotFound)
assert.Nil(t, res)
```

**No verificar el resultado en el camino feliz:**
```go
// ❌ Incompleto: no verifica los valores del resultado
assert.NoError(t, err)

// ✅ Completo: verifica que el resultado tiene los valores esperados
assert.NoError(t, err)
assert.NotNil(t, res)
assert.Equal(t, "valor esperado", res.Campo)
```

---

## 9. Aserciones sin librerías de terceros

Aunque el proyecto usa `github.com/stretchr/testify/assert` por
conveniencia, es importante entender cómo escribir aserciones con la
librería estándar de Go para casos donde testify no esté disponible.

### Equivalencias

| testify | Estándar |
|---------|----------|
| `assert.NoError(t, err)` | `if err != nil { t.Fatalf("error inesperado: %v", err) }` |
| `assert.Error(t, err)` | `if err == nil { t.Fatal("se esperaba un error") }` |
| `assert.ErrorIs(t, err, target)` | `if !errors.Is(err, target) { t.Fatalf("...") }` |
| `assert.Equal(t, a, b)` | `if a != b { t.Fatalf("esperado %v, obtenido %v", a, b) }` |
| `assert.Nil(t, v)` | `if v != nil { t.Fatalf("se esperaba nil, obtenido %v", v) }` |
| `assert.NotNil(t, v)` | `if v == nil { t.Fatal("se esperaba un valor no nil") }` |

### Diferencia entre `t.Error` y `t.Fatal`

- `t.Error` / `assert.Error`: marca el test como fallido pero continúa
  ejecutando el resto del test.
- `t.Fatal` / `require.NoError`: marca el test como fallido y detiene
  la ejecución inmediatamente.

**Regla:** usa `t.Fatal` (o `require`) cuando el resto del test no tiene
sentido si la aserción falla. Por ejemplo, si el resultado es `nil` y
el siguiente paso intenta acceder a un campo, usa `require.NotNil`.

```go
// ✅ Correcto: si res es nil, el acceso a res.Name causaría panic
require.NotNil(t, res)
assert.Equal(t, "esperado", res.Name)

// ❌ Incorrecto: si res es nil, el test paniquea en lugar de fallar limpiamente
assert.NotNil(t, res)
assert.Equal(t, "esperado", res.Name) // panic si res es nil
```

---

## 10. Setup y teardown

### Setup compartido entre subtests

Cuando varios subtests necesitan la misma configuración inicial, crea
el estado compartido antes del primer `t.Run`:

```go
func TestService_Update(t *testing.T) {
    // Setup compartido: se ejecuta una vez para todos los subtests
    repo := &mockRepo{}
    service := NewService(repo)
    userId := "id-de-prueba-valido"

    t.Run("Caso A", func(t *testing.T) { ... })
    t.Run("Caso B", func(t *testing.T) { ... })
}
```

**Importante:** el estado compartido (como `repo`) es mutable. Cada
subtest debe redefinir los campos `Fn` que necesita. No asumas que el
estado del subtest anterior se mantiene.

### Teardown con `t.Cleanup`

Para liberar recursos al final de un test (conexiones, archivos
temporales), usa `t.Cleanup` en lugar de `defer`:

```go
func newMemoryDB(t *testing.T) *gorm.DB {
    t.Helper()
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    require.NoError(t, err)

    // t.Cleanup se ejecuta al final del test, incluso si falla.
    // SQLite en memoria se libera automáticamente al cerrar la conexión.
    sqlDB, _ := db.DB()
    t.Cleanup(func() { sqlDB.Close() })

    require.NoError(t, db.AutoMigrate(&Item{}))
    return db
}
```

**Justificación:** `t.Cleanup` se ejecuta en orden LIFO al final del
test, incluso si el test falla. Es más robusto que `defer` porque
funciona correctamente con subtests.

### Funciones helper de test

Las funciones que crean datos o configuran el entorno para los tests
deben marcarse con `t.Helper()`:

```go
func newMemoryDB(t *testing.T) *gorm.DB {
    t.Helper() // ← hace que los errores apunten al test que llama, no a esta función
    // ...
}
```

**Justificación:** cuando un helper falla, `t.Helper()` hace que el
reporte de error muestre la línea del test que llamó al helper, no la
línea dentro del helper. Esto facilita enormemente la depuración.

---

## 11. Tests de integración vs tests unitarios

### Tests unitarios

Prueban una sola unidad de código (función, método) en completo
aislamiento. Todas las dependencias son mocks.

**Características:**
- No requieren infraestructura externa.
- Se ejecutan en milisegundos.
- Son deterministas: siempre producen el mismo resultado.
- Prueban la lógica de negocio y el manejo de errores.

**Cuándo escribirlos:** para servicios, handlers y cualquier código con
lógica de negocio.

### Tests de integración

Prueban la interacción entre el código y una base de datos real
(SQLite en memoria). Se ejecutan las operaciones reales y se verifica
el estado de la base de datos.

**Características:**
- Usan SQLite en memoria: no requieren servidor externo.
- Son más lentos que los unitarios pero más rápidos que tests con DB real.
- Prueban las queries reales, no una representación de ellas.
- Permiten verificar comportamientos como soft delete y restricciones.

**Cuándo escribirlos:** para repositorios que generan SQL, para
verificar que las operaciones persisten correctamente.

### Diferencia en la estructura

**Test unitario del servicio** (mock del repositorio):
```go
// El repositorio es un mock: controlamos exactamente qué devuelve
repo.GetByIDFn = func(id string) (*SomeType, error) {
    return &SomeType{ID: id}, nil
}
res, err := service.GetByID("id")
```

**Test de integración del repositorio** (SQLite en memoria):
```go
// Se usa una DB real en memoria: se insertan datos y se verifica el resultado
db := newMemoryDB(t)
repo := NewRepository(db)
db.Create(&Item{ID: "id", Name: "nombre"})
res, err := repo.GetByID("id")
```

### Convención para tests de integración del repositorio

Cada subtest crea su propia base de datos SQLite en memoria para
garantizar aislamiento total. El helper `newMemoryDB` migra el esquema
automáticamente:

```go
func newMemoryDB(t *testing.T) *gorm.DB {
    t.Helper()
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    require.NoError(t, err)
    require.NoError(t, db.AutoMigrate(&Item{}))
    return db
}

func TestRepository_GetByID(t *testing.T) {
    t.Run("Debe retornar ErrRecordNotFound cuando el ID no existe", func(t *testing.T) {
        // GIVEN: base de datos vacía
        db := newMemoryDB(t)
        repo := NewRepository(db)

        // WHEN
        result, err := repo.GetByID("id-inexistente")

        // THEN
        assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
        assert.Nil(t, result)
    })

    t.Run("Debe retornar el elemento cuando el ID existe", func(t *testing.T) {
        // GIVEN: elemento insertado directamente en la DB
        db := newMemoryDB(t)
        repo := NewRepository(db)
        require.NoError(t, db.Create(&Item{ID: "id-valido", Name: "Elemento"}).Error)

        // WHEN
        result, err := repo.GetByID("id-valido")

        // THEN
        require.NoError(t, err)
        require.NotNil(t, result)
        assert.Equal(t, "id-valido", result.ID)
    })
}
```

**Ventajas de SQLite en memoria sobre mocks de driver SQL:**
- No hay que declarar queries esperadas: el test es más legible.
- Se prueban las queries reales, no una representación de ellas.
- Permite verificar el estado de la DB después de la operación.
- Permite probar comportamientos como soft delete y restricciones de unicidad.

---

## 12. Qué siempre, qué a veces, qué nunca probar

### Siempre probar

- **Cada rama de error** de una función: si hay tres formas de que
  falle, hay tres tests de error.
- **La transformación de errores:** si la función convierte
  `gorm.ErrRecordNotFound` en `ErrNotFound`, prueba que la
  transformación ocurre correctamente.
- **El camino feliz:** al menos un test donde todo funciona
  correctamente y el resultado tiene los valores esperados.
- **Los efectos secundarios observables:** si la función debe llamar
  a un método del repositorio con ciertos argumentos, verifica que
  lo hace.

### A veces probar (según complejidad)

- **Casos límite de validación:** si la función valida longitudes,
  formatos o rangos, prueba los valores en el límite.
- **Comportamiento con datos vacíos:** listas vacías, strings vacíos,
  valores cero.
- **Idempotencia:** si una operación debe producir el mismo resultado
  al ejecutarse múltiples veces.

### Nunca probar

- **Detalles de implementación privados:** variables privadas, funciones
  no exportadas que no tienen efecto observable desde fuera.
- **El comportamiento de dependencias externas:** no pruebes que GORM
  genera el SQL correcto en un test unitario del servicio; eso es
  responsabilidad del test del repositorio.
- **Código generado automáticamente:** migraciones, código generado por
  herramientas.
- **Configuración de infraestructura:** que el servidor arranca, que
  la base de datos acepta conexiones.

---

## 13. Anti-patrones a evitar

### 1. El test que siempre pasa

```go
// ❌ Este test siempre pasa porque no hay aserciones
func TestService_GetByID(t *testing.T) {
    repo := &mockRepo{}
    service := NewService(repo)
    repo.GetByIDFn = func(id string) (*SomeType, error) { return nil, nil }
    service.GetByID("id") // sin aserciones
}
```

**Por qué es dañino:** da falsa confianza. El código puede estar roto
y el test seguirá pasando.

### 2. El test que prueba demasiado

```go
// ❌ Un solo test que cubre múltiples escenarios
func TestService_Update(t *testing.T) {
    // Caso 1
    repo.GetByIDFn = func(id string) (*SomeType, error) { return nil, ErrNotFound }
    _, err1 := service.Update("id", req)
    assert.Error(t, err1)

    // Caso 2 (depende del estado del caso 1)
    repo.GetByIDFn = func(id string) (*SomeType, error) { return &SomeType{}, nil }
    _, err2 := service.Update("id", req)
    assert.NoError(t, err2)
}
```

**Por qué es dañino:** cuando falla, no sabes cuál de los casos causó
el fallo. Usa `t.Run` para separar los escenarios.

### 3. Mocks que no verifican sus argumentos

```go
// ❌ El mock acepta cualquier argumento sin verificar
repo.UpdateFn = func(u *SomeType) error {
    return nil // no verifica que u tiene los valores correctos
}
```

**Por qué es dañino:** el test pasa aunque el servicio llame a
`Update` con datos incorrectos. Agrega aserciones dentro del mock
cuando los argumentos importan.

### 4. Aserciones sobre el mensaje de error en lugar del tipo

```go
// ❌ Frágil: falla si se cambia el mensaje de error
assert.EqualError(t, err, "Usuario no encontrado")

// ✅ Robusto: verifica el tipo de error, no el mensaje
assert.ErrorIs(t, err, ErrNotFound)
```

**Por qué es dañino:** los mensajes de error cambian (traducciones,
mejoras de UX). El tipo de error es parte del contrato de la API.

### 5. Tests que dependen del orden de ejecución

```go
// ❌ El test B depende del estado que dejó el test A
var sharedRepo = &mockRepo{}

func TestA(t *testing.T) {
    sharedRepo.GetByIDFn = func(id string) (*SomeType, error) { ... }
    // ...
}

func TestB(t *testing.T) {
    // Asume que GetByIDFn tiene el valor que dejó TestA
    // ...
}
```

**Por qué es dañino:** Go no garantiza el orden de ejecución de las
funciones `TestXxx`. Cada test debe ser completamente independiente.

### 6. Ignorar el resultado de `mock.ExpectationsWereMet()`

```go
// ❌ No verifica que todas las queries esperadas se ejecutaron
mock.ExpectQuery(`SELECT \* FROM "items"`).WillReturnRows(...)
result, err := repo.GetByID("id")
assert.NoError(t, err)
// Falta: assert.NoError(t, mock.ExpectationsWereMet())
```

**Por qué es dañino:** el test puede pasar aunque el repositorio no
haya ejecutado la query esperada.

---

## 14. Ejemplo completo de archivo de test unitario

El siguiente ejemplo muestra un archivo de test unitario completo que
sigue todas las convenciones de esta guía. Está diseñado para ser
copiado como plantilla al crear tests para un nuevo servicio.

```go
package example_test

// Los tests del servicio verifican la lógica de negocio de cada método.
// Se usa un mock manual del repositorio para aislar completamente el
// servicio de la base de datos.
//
// Estrategia:
//   - mockRepo implementa Repository con campos de tipo func.
//   - Cada subtest redefine solo los campos que necesita.
//   - Se sigue el patrón GIVEN / WHEN / THEN en cada caso.

import (
    "errors"
    "testing"

    "github.com/stretchr/testify/assert"
    "gorm.io/gorm"
)

// ─────────────────────────────────────────────────────────────────────────────
// Mock del repositorio
// ─────────────────────────────────────────────────────────────────────────────

type mockRepo struct {
    GetByIDFn func(id string) (*Item, error)
    CreateFn  func(item *Item) error
    UpdateFn  func(item *Item) error
    DeleteFn  func(id string) error
}

func (m *mockRepo) GetByID(id string) (*Item, error) { return m.GetByIDFn(id) }
func (m *mockRepo) Create(item *Item) error          { return m.CreateFn(item) }
func (m *mockRepo) Update(item *Item) error          { return m.UpdateFn(item) }
func (m *mockRepo) Delete(id string) error           { return m.DeleteFn(id) }

// ─────────────────────────────────────────────────────────────────────────────
// GetByID
// ─────────────────────────────────────────────────────────────────────────────

func TestService_GetByID(t *testing.T) {
    repo := &mockRepo{}
    service := NewService(repo)
    itemId := "id-de-prueba"

    // ----------------------------------------------------------------
    // Caso 1: elemento no existe → ErrNotFound
    // ----------------------------------------------------------------
    t.Run("Debe retornar ErrNotFound cuando el elemento no existe", func(t *testing.T) {
        // GIVEN: el repositorio responde con ErrRecordNotFound
        repo.GetByIDFn = func(id string) (*Item, error) {
            return nil, gorm.ErrRecordNotFound
        }

        // WHEN
        res, err := service.GetByID(itemId)

        // THEN: el error de GORM se transforma en el error de dominio correcto
        assert.ErrorIs(t, err, ErrNotFound)
        assert.Nil(t, res)
    })

    // ----------------------------------------------------------------
    // Caso 2: error genérico del repositorio → se propaga tal cual
    // ----------------------------------------------------------------
    t.Run("Debe propagar el error genérico cuando el repositorio falla", func(t *testing.T) {
        // GIVEN: error de conexión que no es ErrRecordNotFound
        dbErr := errors.New("connection refused")
        repo.GetByIDFn = func(id string) (*Item, error) {
            return nil, dbErr
        }

        // WHEN
        res, err := service.GetByID(itemId)

        // THEN: el error llega sin transformar y NO es ErrNotFound
        assert.ErrorIs(t, err, dbErr)
        assert.NotErrorIs(t, err, ErrNotFound)
        assert.Nil(t, res)
    })

    // ----------------------------------------------------------------
    // Caso 3: elemento encontrado exitosamente
    // ----------------------------------------------------------------
    t.Run("Debe retornar el elemento cuando existe", func(t *testing.T) {
        // GIVEN: el repositorio devuelve un elemento válido
        expected := &Item{ID: itemId, Name: "Elemento de prueba"}
        repo.GetByIDFn = func(id string) (*Item, error) {
            return expected, nil
        }

        // WHEN
        res, err := service.GetByID(itemId)

        // THEN
        assert.NoError(t, err)
        assert.Equal(t, expected, res)
    })
}

// ─────────────────────────────────────────────────────────────────────────────
// Update
// ─────────────────────────────────────────────────────────────────────────────

func TestService_Update(t *testing.T) {
    repo := &mockRepo{}
    service := NewService(repo)
    itemId := "id-de-prueba"

    // ----------------------------------------------------------------
    // Caso 1: elemento no existe → ErrNotFound
    // ----------------------------------------------------------------
    t.Run("Debe retornar ErrNotFound cuando el elemento no existe", func(t *testing.T) {
        // GIVEN
        repo.GetByIDFn = func(id string) (*Item, error) {
            return nil, gorm.ErrRecordNotFound
        }

        // WHEN
        res, err := service.Update(itemId, UpdateRequest{Name: "Nuevo"})

        // THEN
        assert.ErrorIs(t, err, ErrNotFound)
        assert.Nil(t, res)
    })

    // ----------------------------------------------------------------
    // Caso 2: repo.Update falla → error propagado
    // ----------------------------------------------------------------
    t.Run("Debe propagar el error del repositorio cuando Update falla", func(t *testing.T) {
        // GIVEN: el elemento existe pero la actualización falla
        repo.GetByIDFn = func(id string) (*Item, error) {
            return &Item{ID: itemId, Name: "Original"}, nil
        }
        repo.UpdateFn = func(item *Item) error {
            return errors.New("constraint violation")
        }

        // WHEN
        res, err := service.Update(itemId, UpdateRequest{Name: "Nuevo"})

        // THEN
        assert.EqualError(t, err, "constraint violation")
        assert.Nil(t, res)
    })

    // ----------------------------------------------------------------
    // Caso 3: actualización exitosa → elemento con valores nuevos
    // ----------------------------------------------------------------
    t.Run("Debe actualizar el elemento correctamente cuando todo es válido", func(t *testing.T) {
        // GIVEN: el elemento existe y la actualización tiene éxito
        nuevoNombre := "Nombre Actualizado"
        repo.GetByIDFn = func(id string) (*Item, error) {
            return &Item{ID: itemId, Name: "Original"}, nil
        }
        repo.UpdateFn = func(item *Item) error {
            // Verificamos que el objeto llegó con los valores correctos
            assert.Equal(t, nuevoNombre, item.Name)
            return nil
        }

        // WHEN
        res, err := service.Update(itemId, UpdateRequest{Name: nuevoNombre})

        // THEN
        assert.NoError(t, err)
        assert.NotNil(t, res)
        assert.Equal(t, nuevoNombre, res.Name)
        assert.Equal(t, itemId, res.ID) // el ID no debe cambiar
    })
}

// ─────────────────────────────────────────────────────────────────────────────
// Delete
// ─────────────────────────────────────────────────────────────────────────────

func TestService_Delete(t *testing.T) {
    repo := &mockRepo{}
    service := NewService(repo)
    itemId := "id-de-prueba"

    // ----------------------------------------------------------------
    // Caso 1: elemento no existe → ErrNotFound
    // ----------------------------------------------------------------
    t.Run("Debe retornar ErrNotFound cuando el elemento no existe", func(t *testing.T) {
        // GIVEN
        repo.DeleteFn = func(id string) error {
            return gorm.ErrRecordNotFound
        }

        // WHEN
        err := service.Delete(itemId)

        // THEN
        assert.ErrorIs(t, err, ErrNotFound)
    })

    // ----------------------------------------------------------------
    // Caso 2: error genérico del repositorio → se propaga
    // ----------------------------------------------------------------
    t.Run("Debe propagar el error genérico cuando el repositorio falla", func(t *testing.T) {
        // GIVEN
        repo.DeleteFn = func(id string) error {
            return errors.New("db error")
        }

        // WHEN
        err := service.Delete(itemId)

        // THEN
        assert.EqualError(t, err, "db error")
    })

    // ----------------------------------------------------------------
    // Caso 3: eliminación exitosa
    // ----------------------------------------------------------------
    t.Run("Debe eliminar el elemento sin error cuando existe", func(t *testing.T) {
        // GIVEN
        repo.DeleteFn = func(id string) error {
            return nil
        }

        // WHEN
        err := service.Delete(itemId)

        // THEN
        assert.NoError(t, err)
    })
}
```

---

## Referencia rápida de aserciones

```go
import "github.com/stretchr/testify/assert"
import "github.com/stretchr/testify/require"
```

| Aserción | Cuándo usarla |
|----------|---------------|
| `assert.NoError(t, err)` | El resultado no debe tener error |
| `assert.Error(t, err)` | Debe existir algún error (usar solo cuando el tipo no importa) |
| `assert.ErrorIs(t, err, target)` | El error es (o envuelve) `target` |
| `assert.EqualError(t, err, "msg")` | El mensaje del error es exactamente `"msg"` |
| `assert.Nil(t, val)` | El valor es `nil` |
| `assert.NotNil(t, val)` | El valor no es `nil` |
| `assert.Equal(t, expected, actual)` | Igualdad profunda (structs, slices, etc.) |
| `assert.NotEqual(t, a, b)` | Los valores son distintos |
| `assert.Len(t, slice, n)` | El slice tiene exactamente `n` elementos |
| `assert.Empty(t, val)` | El valor está vacío (nil, "", [], {}) |
| `assert.True(t, cond)` | La condición es verdadera |
| `assert.False(t, cond)` | La condición es falsa |
| `assert.Contains(t, s, substr)` | El string contiene el substring |
| `require.NoError(t, err)` | Como `assert.NoError` pero detiene el test si falla |
| `require.NotNil(t, val)` | Como `assert.NotNil` pero detiene el test si falla |

### Cuándo usar `require` en lugar de `assert`

Usa `require` cuando el resto del test no tiene sentido si la aserción
falla. El caso más común es verificar que un resultado no es `nil` antes
de acceder a sus campos:

```go
res, err := service.GetByID("id")
require.NoError(t, err)    // si hay error, no tiene sentido continuar
require.NotNil(t, res)     // si res es nil, el siguiente acceso causaría panic
assert.Equal(t, "esperado", res.Name)
```

---

## Comandos de ejecución

```bash
# Todos los tests del proyecto
go test ./...

# Tests de un módulo específico
go test ./internal/modules/example/...

# Con salida detallada (ver nombre de cada subtest)
go test -v ./internal/modules/example/...

# Solo los tests que coincidan con un patrón
go test -v -run "TestService_Update" ./internal/modules/example/...

# Con reporte de cobertura en consola
go test -cover ./internal/modules/example/...

# Generar reporte HTML de cobertura
go test -coverprofile=coverage.out ./internal/modules/example/...
go tool cover -html=coverage.out
```
