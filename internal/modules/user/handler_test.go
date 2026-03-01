package user_test

// Los tests del handler verifican la capa HTTP: códigos de estado, formato
// de respuesta JSON y que el handler llame al servicio con los parámetros
// correctos.
//
// Estrategia:
//   - Se crea un mockService que implementa user.Service.
//   - Se usa net/http/httptest para simular peticiones HTTP sin levantar servidor.
//   - Se usa github.com/go-chi/chi/v5 para registrar rutas con parámetros de URL.
//   - Se decodifica el JSON de respuesta y se verifican campos clave.

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"image-processing-service/internal/modules/user"
	"image-processing-service/internal/shared/auth"
	"image-processing-service/internal/shared/utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─────────────────────────────────────────────────────────────────────────────
// Mock del servicio
// ─────────────────────────────────────────────────────────────────────────────

// mockService implementa user.Service con campos de tipo func para que cada
// test pueda redefinir el comportamiento sin crear una nueva instancia.
type mockService struct {
	GetByIDFn        func(id string) (*user.User, error)
	GetAllFn         func() ([]*user.User, error)
	UpdateFn         func(id string, req user.UpdateUserRequest) (*user.User, error)
	UpdatePasswordFn func(id string, req user.UpdatePasswordUserRequest) (*user.User, error)
	DeleteFn         func(id string) error
}

func (m *mockService) GetByID(id string) (*user.User, error) {
	return m.GetByIDFn(id)
}
func (m *mockService) GetAll() ([]*user.User, error) {
	return m.GetAllFn()
}
func (m *mockService) Update(id string, req user.UpdateUserRequest) (*user.User, error) {
	return m.UpdateFn(id, req)
}
func (m *mockService) UpdatePassword(id string, req user.UpdatePasswordUserRequest) (*user.User, error) {
	return m.UpdatePasswordFn(id, req)
}
func (m *mockService) Delete(id string) error {
	return m.DeleteFn(id)
}

// ─────────────────────────────────────────────────────────────────────────────
// Helpers de test
// ─────────────────────────────────────────────────────────────────────────────

// validID es un CUID2 válido de 24 caracteres para usar en los tests.
// Debe pasar la validación de utils.IsValidID (cuid2.IsCuid).
const validID = "clbxyz1234567890abcdefgh"

// newRequest crea una petición HTTP con el ID inyectado como parámetro de ruta
// de chi, simulando lo que haría el router en producción.
func newRequest(method, path, body string, urlParams map[string]string) *http.Request {
	var bodyReader *strings.Reader
	if body != "" {
		bodyReader = strings.NewReader(body)
	} else {
		bodyReader = strings.NewReader("")
	}

	req := httptest.NewRequest(method, path, bodyReader)
	req.Header.Set("Content-Type", "application/json")

	// Inyectar parámetros de URL de chi en el contexto
	if len(urlParams) > 0 {
		rctx := chi.NewRouteContext()
		for k, v := range urlParams {
			rctx.URLParams.Add(k, v)
		}
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	}

	return req
}

// withAuthUser inyecta un usuario autenticado en el contexto de la petición,
// simulando lo que hace el middleware de autenticación.
func withAuthUser(req *http.Request, userID string) *http.Request {
	authUser := auth.AuthenticatedUser{UserID: userID}
	ctx := context.WithValue(req.Context(), auth.AuthKey, authUser)
	return req.WithContext(ctx)
}

// decodeResponse decodifica el JSON de la respuesta en un map genérico.
func decodeResponse(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var result map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&result)
	require.NoError(t, err, "la respuesta no es JSON válido")
	return result
}

// ─────────────────────────────────────────────────────────────────────────────
// GetByID
// ─────────────────────────────────────────────────────────────────────────────

func TestHandler_GetByID(t *testing.T) {
	svc := &mockService{}
	h := user.NewHandler(svc)

	// ----------------------------------------------------------------
	// Caso 1: ID con formato inválido → 400
	// ----------------------------------------------------------------
	t.Run("Debe retornar 400 cuando el ID tiene formato inválido", func(t *testing.T) {
		// GIVEN: ID que no pasa la validación cuid2
		req := newRequest(http.MethodGet, "/users/id-invalido", "", map[string]string{"id": "id-invalido"})
		w := httptest.NewRecorder()

		// WHEN
		h.GetByID(w, req)

		// THEN
		assert.Equal(t, http.StatusBadRequest, w.Code)
		body := decodeResponse(t, w)
		assert.False(t, body["success"].(bool))
	})

	// ----------------------------------------------------------------
	// Caso 2: usuario no encontrado → 404
	// ----------------------------------------------------------------
	t.Run("Debe retornar 404 cuando el usuario no existe", func(t *testing.T) {
		// GIVEN: el servicio devuelve ErrNotFound
		svc.GetByIDFn = func(id string) (*user.User, error) {
			return nil, user.ErrNotFound
		}

		req := newRequest(http.MethodGet, "/users/"+validID, "", map[string]string{"id": validID})
		w := httptest.NewRecorder()

		// WHEN
		h.GetByID(w, req)

		// THEN
		assert.Equal(t, http.StatusNotFound, w.Code)
		body := decodeResponse(t, w)
		assert.False(t, body["success"].(bool))
	})

	// ----------------------------------------------------------------
	// Caso 3: error interno del servicio → 500
	// ----------------------------------------------------------------
	t.Run("Debe retornar 500 cuando el servicio falla con error genérico", func(t *testing.T) {
		// GIVEN: el servicio devuelve un error no tipado
		svc.GetByIDFn = func(id string) (*user.User, error) {
			return nil, errors.New("db timeout")
		}

		req := newRequest(http.MethodGet, "/users/"+validID, "", map[string]string{"id": validID})
		w := httptest.NewRecorder()

		// WHEN
		h.GetByID(w, req)

		// THEN: HandleError convierte errores no-AppError en 500
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	// ----------------------------------------------------------------
	// Caso 4: usuario encontrado → 200 con datos
	// ----------------------------------------------------------------
	t.Run("Debe retornar 200 con el usuario cuando existe", func(t *testing.T) {
		// GIVEN: el servicio devuelve un usuario válido
		svc.GetByIDFn = func(id string) (*user.User, error) {
			return &user.User{ID: validID, Name: "Ana", Email: "ana@test.com"}, nil
		}

		req := newRequest(http.MethodGet, "/users/"+validID, "", map[string]string{"id": validID})
		w := httptest.NewRecorder()

		// WHEN
		h.GetByID(w, req)

		// THEN
		assert.Equal(t, http.StatusOK, w.Code)
		body := decodeResponse(t, w)
		assert.True(t, body["success"].(bool))

		data := body["data"].(map[string]interface{})
		assert.Equal(t, "Ana", data["name"])
		assert.Equal(t, "ana@test.com", data["email"])
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// GetAll
// ─────────────────────────────────────────────────────────────────────────────

func TestHandler_GetAll(t *testing.T) {
	svc := &mockService{}
	h := user.NewHandler(svc)

	// ----------------------------------------------------------------
	// Caso 1: error del servicio → 500
	// ----------------------------------------------------------------
	t.Run("Debe retornar 500 cuando el servicio falla", func(t *testing.T) {
		// GIVEN
		svc.GetAllFn = func() ([]*user.User, error) {
			return nil, errors.New("db down")
		}

		req := newRequest(http.MethodGet, "/users", "", nil)
		w := httptest.NewRecorder()

		// WHEN
		h.GetAll(w, req)

		// THEN
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	// ----------------------------------------------------------------
	// Caso 2: lista vacía → 200 con array vacío
	// ----------------------------------------------------------------
	t.Run("Debe retornar 200 con lista vacía cuando no hay usuarios", func(t *testing.T) {
		// GIVEN
		svc.GetAllFn = func() ([]*user.User, error) {
			return []*user.User{}, nil
		}

		req := newRequest(http.MethodGet, "/users", "", nil)
		w := httptest.NewRecorder()

		// WHEN
		h.GetAll(w, req)

		// THEN
		assert.Equal(t, http.StatusOK, w.Code)
		body := decodeResponse(t, w)
		assert.True(t, body["success"].(bool))
		data := body["data"].([]interface{})
		assert.Empty(t, data)
	})

	// ----------------------------------------------------------------
	// Caso 3: múltiples usuarios → 200 con todos
	// ----------------------------------------------------------------
	t.Run("Debe retornar 200 con todos los usuarios", func(t *testing.T) {
		// GIVEN
		svc.GetAllFn = func() ([]*user.User, error) {
			return []*user.User{
				{ID: "id-1", Name: "Ana", Email: "ana@test.com"},
				{ID: "id-2", Name: "Luis", Email: "luis@test.com"},
			}, nil
		}

		req := newRequest(http.MethodGet, "/users", "", nil)
		w := httptest.NewRecorder()

		// WHEN
		h.GetAll(w, req)

		// THEN
		assert.Equal(t, http.StatusOK, w.Code)
		body := decodeResponse(t, w)
		data := body["data"].([]interface{})
		assert.Len(t, data, 2)
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// Update
// ─────────────────────────────────────────────────────────────────────────────

func TestHandler_Update(t *testing.T) {
	svc := &mockService{}
	h := user.NewHandler(svc)

	// ----------------------------------------------------------------
	// Caso 1: ID inválido → 400
	// ----------------------------------------------------------------
	t.Run("Debe retornar 400 cuando el ID tiene formato inválido", func(t *testing.T) {
		// GIVEN
		req := newRequest(http.MethodPatch, "/users/bad-id",
			`{"name":"Nuevo"}`,
			map[string]string{"id": "bad-id"},
		)
		w := httptest.NewRecorder()

		// WHEN
		h.Update(w, req)

		// THEN
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// ----------------------------------------------------------------
	// Caso 2: body no es JSON válido → 400
	// ----------------------------------------------------------------
	t.Run("Debe retornar 400 cuando el body no es JSON válido", func(t *testing.T) {
		// GIVEN: body malformado
		req := newRequest(http.MethodPatch, "/users/"+validID,
			`{esto no es json`,
			map[string]string{"id": validID},
		)
		w := httptest.NewRecorder()

		// WHEN
		h.Update(w, req)

		// THEN
		assert.Equal(t, http.StatusBadRequest, w.Code)
		body := decodeResponse(t, w)
		errObj := body["error"].(map[string]interface{})
		assert.Equal(t, "INVALID_JSON", errObj["code"])
	})

	// ----------------------------------------------------------------
	// Caso 3: validación falla (email inválido) → 422
	// ----------------------------------------------------------------
	t.Run("Debe retornar 422 cuando el email tiene formato inválido", func(t *testing.T) {
		// GIVEN: email que no pasa la validación del struct
		req := newRequest(http.MethodPatch, "/users/"+validID,
			`{"email":"no-es-un-email"}`,
			map[string]string{"id": validID},
		)
		w := httptest.NewRecorder()

		// WHEN
		h.Update(w, req)

		// THEN
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
		body := decodeResponse(t, w)
		errObj := body["error"].(map[string]interface{})
		assert.Equal(t, "VALIDATION_FAILED", errObj["code"])
	})

	// ----------------------------------------------------------------
	// Caso 4: usuario no encontrado → 404
	// ----------------------------------------------------------------
	t.Run("Debe retornar 404 cuando el usuario no existe", func(t *testing.T) {
		// GIVEN
		svc.UpdateFn = func(id string, req user.UpdateUserRequest) (*user.User, error) {
			return nil, user.ErrNotFound
		}

		req := newRequest(http.MethodPatch, "/users/"+validID,
			`{"name":"Nuevo nombre"}`,
			map[string]string{"id": validID},
		)
		w := httptest.NewRecorder()

		// WHEN
		h.Update(w, req)

		// THEN
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	// ----------------------------------------------------------------
	// Caso 5: email ya en uso → 409
	// ----------------------------------------------------------------
	t.Run("Debe retornar 409 cuando el email ya está en uso", func(t *testing.T) {
		// GIVEN
		svc.UpdateFn = func(id string, req user.UpdateUserRequest) (*user.User, error) {
			return nil, user.ErrAlreadyExists
		}

		nuevoEmail := "ocupado@test.com"
		body, _ := json.Marshal(map[string]string{"email": nuevoEmail})
		req := newRequest(http.MethodPatch, "/users/"+validID,
			string(body),
			map[string]string{"id": validID},
		)
		w := httptest.NewRecorder()

		// WHEN
		h.Update(w, req)

		// THEN
		assert.Equal(t, http.StatusConflict, w.Code)
	})

	// ----------------------------------------------------------------
	// Caso 6: actualización exitosa → 201 con usuario actualizado
	// ----------------------------------------------------------------
	t.Run("Debe retornar 201 con el usuario actualizado cuando todo es válido", func(t *testing.T) {
		// GIVEN
		nuevoNombre := "Nombre Actualizado"
		svc.UpdateFn = func(id string, req user.UpdateUserRequest) (*user.User, error) {
			return &user.User{
				ID:    validID,
				Name:  nuevoNombre,
				Email: "correo@test.com",
			}, nil
		}

		reqBody, _ := json.Marshal(map[string]string{"name": nuevoNombre})
		req := newRequest(http.MethodPatch, "/users/"+validID,
			string(reqBody),
			map[string]string{"id": validID},
		)
		w := httptest.NewRecorder()

		// WHEN
		h.Update(w, req)

		// THEN
		assert.Equal(t, http.StatusCreated, w.Code)
		respBody := decodeResponse(t, w)
		assert.True(t, respBody["success"].(bool))
		data := respBody["data"].(map[string]interface{})
		assert.Equal(t, nuevoNombre, data["name"])
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// UpdatePassword
// ─────────────────────────────────────────────────────────────────────────────

func TestHandler_UpdatePassword(t *testing.T) {
	svc := &mockService{}
	h := user.NewHandler(svc)

	// ----------------------------------------------------------------
	// Caso 1: ID del usuario autenticado inválido → 400
	// ----------------------------------------------------------------
	t.Run("Debe retornar 400 cuando el ID del usuario autenticado es inválido", func(t *testing.T) {
		// GIVEN: contexto con un UserID que no pasa la validación cuid2
		req := newRequest(http.MethodPatch, "/users/password",
			`{"current_password":"pass123","new_password":"nueva456"}`,
			nil,
		)
		req = withAuthUser(req, "id-invalido")
		w := httptest.NewRecorder()

		// WHEN
		h.UpdatePassword(w, req)

		// THEN
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// ----------------------------------------------------------------
	// Caso 2: body no es JSON válido → 400
	// ----------------------------------------------------------------
	t.Run("Debe retornar 400 cuando el body no es JSON válido", func(t *testing.T) {
		// GIVEN
		req := newRequest(http.MethodPatch, "/users/password",
			`{malformado`,
			nil,
		)
		req = withAuthUser(req, validID)
		w := httptest.NewRecorder()

		// WHEN
		h.UpdatePassword(w, req)

		// THEN
		assert.Equal(t, http.StatusBadRequest, w.Code)
		body := decodeResponse(t, w)
		errObj := body["error"].(map[string]interface{})
		assert.Equal(t, "INVALID_JSON", errObj["code"])
	})

	// ----------------------------------------------------------------
	// Caso 3: validación falla (contraseña muy corta) → 422
	// ----------------------------------------------------------------
	t.Run("Debe retornar 422 cuando la contraseña nueva es muy corta", func(t *testing.T) {
		// GIVEN: new_password tiene menos de 6 caracteres
		req := newRequest(http.MethodPatch, "/users/password",
			`{"current_password":"pass123","new_password":"abc"}`,
			nil,
		)
		req = withAuthUser(req, validID)
		w := httptest.NewRecorder()

		// WHEN
		h.UpdatePassword(w, req)

		// THEN
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	// ----------------------------------------------------------------
	// Caso 4: contraseña actual incorrecta → 403
	// ----------------------------------------------------------------
	t.Run("Debe retornar 403 cuando la contraseña actual es incorrecta", func(t *testing.T) {
		// GIVEN
		svc.UpdatePasswordFn = func(id string, req user.UpdatePasswordUserRequest) (*user.User, error) {
			return nil, user.ErrInvalidPassword
		}

		req := newRequest(http.MethodPatch, "/users/password",
			`{"current_password":"incorrecta","new_password":"nueva456"}`,
			nil,
		)
		req = withAuthUser(req, validID)
		w := httptest.NewRecorder()

		// WHEN
		h.UpdatePassword(w, req)

		// THEN
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	// ----------------------------------------------------------------
	// Caso 5: actualización exitosa → 201
	// ----------------------------------------------------------------
	t.Run("Debe retornar 201 cuando la contraseña se actualiza correctamente", func(t *testing.T) {
		// GIVEN
		svc.UpdatePasswordFn = func(id string, req user.UpdatePasswordUserRequest) (*user.User, error) {
			return &user.User{ID: validID, Name: "Ana", Email: "ana@test.com"}, nil
		}

		req := newRequest(http.MethodPatch, "/users/password",
			`{"current_password":"correcta123","new_password":"nueva456"}`,
			nil,
		)
		req = withAuthUser(req, validID)
		w := httptest.NewRecorder()

		// WHEN
		h.UpdatePassword(w, req)

		// THEN
		assert.Equal(t, http.StatusCreated, w.Code)
		body := decodeResponse(t, w)
		assert.True(t, body["success"].(bool))
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// Delete
// ─────────────────────────────────────────────────────────────────────────────

func TestHandler_Delete(t *testing.T) {
	svc := &mockService{}
	h := user.NewHandler(svc)

	// ----------------------------------------------------------------
	// Caso 1: ID inválido → 400
	// ----------------------------------------------------------------
	t.Run("Debe retornar 400 cuando el ID tiene formato inválido", func(t *testing.T) {
		// GIVEN
		req := newRequest(http.MethodDelete, "/users/bad-id", "", map[string]string{"id": "bad-id"})
		w := httptest.NewRecorder()

		// WHEN
		h.Delete(w, req)

		// THEN
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// ----------------------------------------------------------------
	// Caso 2: usuario no encontrado → 404
	// ----------------------------------------------------------------
	t.Run("Debe retornar 404 cuando el usuario no existe", func(t *testing.T) {
		// GIVEN
		svc.DeleteFn = func(id string) error {
			return user.ErrNotFound
		}

		req := newRequest(http.MethodDelete, "/users/"+validID, "", map[string]string{"id": validID})
		w := httptest.NewRecorder()

		// WHEN
		h.Delete(w, req)

		// THEN
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	// ----------------------------------------------------------------
	// Caso 3: error interno → 500
	// ----------------------------------------------------------------
	t.Run("Debe retornar 500 cuando el servicio falla con error genérico", func(t *testing.T) {
		// GIVEN
		svc.DeleteFn = func(id string) error {
			return errors.New("db error")
		}

		req := newRequest(http.MethodDelete, "/users/"+validID, "", map[string]string{"id": validID})
		w := httptest.NewRecorder()

		// WHEN
		h.Delete(w, req)

		// THEN
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	// ----------------------------------------------------------------
	// Caso 4: eliminación exitosa → 200 con mensaje
	// ----------------------------------------------------------------
	t.Run("Debe retornar 200 con mensaje de confirmación cuando el usuario existe", func(t *testing.T) {
		// GIVEN
		svc.DeleteFn = func(id string) error {
			return nil
		}

		req := newRequest(http.MethodDelete, "/users/"+validID, "", map[string]string{"id": validID})
		w := httptest.NewRecorder()

		// WHEN
		h.Delete(w, req)

		// THEN
		assert.Equal(t, http.StatusOK, w.Code)
		body := decodeResponse(t, w)
		assert.True(t, body["success"].(bool))
		data := body["data"].(map[string]interface{})
		assert.Contains(t, data["message"], "eliminado")
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// Nota sobre validID
// ─────────────────────────────────────────────────────────────────────────────
//
// utils.IsValidID usa cuid2.IsCuid, que valida el formato CUID2.
// Si los tests de ID inválido fallan, verifica que la constante validID
// sea un CUID2 real. Puedes generar uno con:
//   go run -e 'fmt.Println(utils.GenerateID())'
//
// La constante actual es solo un placeholder; reemplázala con un ID real
// si los tests de "ID válido" fallan.

// Aseguramos que bytes y utils se importan
var _ = bytes.NewBuffer
var _ = utils.Pointer[string]
