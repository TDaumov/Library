package handler

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"library_vebservice/dto"
	customMid "library_vebservice/middleware"
	"library_vebservice/models"
	"library_vebservice/service"
	"net/http"
	"time"
)

type Handler struct {
	authService service.AuthService
	bookService service.BookService
	logger      zerolog.Logger
	validate    *validator.Validate
	jwtSecret   string
}

func NewHandler(a service.AuthService, b service.BookService, log zerolog.Logger, jwtSecret string) *Handler {
	return &Handler{
		authService: a,
		bookService: b,
		logger:      log,
		validate:    validator.New(),
		jwtSecret:   jwtSecret,
	}
}

func (h *Handler) RegisterRoutes(e *echo.Echo) {
	v1 := e.Group("/api/v1")

	// Публичные маршруты
	auth := v1.Group("/auth")
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)

	// Защищённые маршруты
	books := v1.Group("/books")
	books.Use(customMid.JWTMiddlewareFunc())

	books.GET("", h.ListBooks)
	books.GET("/:id", h.GetBook)
	books.POST("", h.CreateBook)
	books.PUT("/:id", h.UpdateBook)
	books.DELETE("/:id", h.DeleteBook)
	books.POST("/:id/borrow", h.BorrowBook)
}

// ------------------- Implementations -------------------

func (h *Handler) Register(c echo.Context) error {
	var dto dto.RegisterDTO
	if err := c.Bind(&dto); err != nil {
		h.logger.Error().Err(err).Msg("Failed to bind RegisterDTO")
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if err := h.validate.Struct(dto); err != nil {
		h.logger.Error().Err(err).Msg("Validation failed for RegisterDTO")
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	user, err := h.authService.Register(context.Background(), dto.Email, dto.Password, dto.Role)
	if err != nil {
		h.logger.Error().Err(err).Str("email", dto.Email).Msg("Registration failed")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	h.logger.Info().Str("email", dto.Email).Str("role", dto.Role).Msg("User registered successfully")
	return c.JSON(http.StatusCreated, map[string]interface{}{"id": user.ID, "email": user.Email, "role": user.Role})
}

func (h *Handler) Login(c echo.Context) error {
	var l dto.LoginDTO
	if err := c.Bind(&l); err != nil {
		h.logger.Error().Err(err).Msg("Failed to bind LoginDTO")
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if err := h.validate.Struct(l); err != nil {
		h.logger.Error().Err(err).Msg("Validation failed for LoginDTO")
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	token, err := h.authService.Login(context.Background(), l.Email, l.Password, h.jwtSecret)
	if err != nil {
		h.logger.Warn().Str("email", l.Email).Err(err).Msg("Login failed")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
	h.logger.Info().Str("email", l.Email).Msg("User logged in successfully")
	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) ListBooks(c echo.Context) error {
	list, err := h.bookService.List(context.Background())
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to list books")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	h.logger.Info().Int("count", len(list)).Msg("Listed books successfully")
	return c.JSON(http.StatusOK, list)
}

func (h *Handler) GetBook(c echo.Context) error {
	id := c.Param("id")
	b, err := h.bookService.GetByID(context.Background(), id)
	if err != nil {
		h.logger.Warn().Str("book_id", id).Msg("Book not found")
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	h.logger.Info().Str("book_id", id).Msg("Fetched book successfully")
	return c.JSON(http.StatusOK, b)
}

func (h *Handler) CreateBook(c echo.Context) error {
	var d dto.CreateBookDTO
	if err := c.Bind(&d); err != nil {
		h.logger.Error().Err(err).Msg("Failed to bind CreateBookDTO")
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if err := h.validate.Struct(d); err != nil {
		h.logger.Error().Err(err).Msg("Validation failed for CreateBookDTO")
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	b := &models.Book{
		Title: d.Title, Author: d.Author, Description: d.Description, PublishedAt: d.PublishedAt,
	}
	if err := h.bookService.Create(context.Background(), b); err != nil {
		h.logger.Error().Err(err).Str("title", d.Title).Msg("Failed to create book")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	h.logger.Info().Str("book_id", b.ID.String()).Str("title", b.Title).Msg("Book created successfully")
	return c.JSON(http.StatusCreated, b)
}

func (h *Handler) UpdateBook(c echo.Context) error {
	id := c.Param("id")
	var u dto.UpdateBookDTO
	if err := c.Bind(&u); err != nil {
		h.logger.Error().Err(err).Msg("Failed to bind UpdateBookDTO")
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	book, err := h.bookService.GetByID(context.Background(), id)
	if err != nil {
		h.logger.Warn().Str("book_id", id).Msg("Book not found for update")
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	if u.Title != nil {
		book.Title = *u.Title
	}
	if u.Author != nil {
		book.Author = *u.Author
	}
	if u.Description != nil {
		book.Description = *u.Description
	}
	if err := h.bookService.Update(context.Background(), book); err != nil {
		h.logger.Error().Err(err).Str("book_id", id).Msg("Failed to update book")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	h.logger.Info().Str("book_id", id).Msg("Book updated successfully")
	return c.JSON(http.StatusOK, book)
}

func (h *Handler) DeleteBook(c echo.Context) error {
	id := c.Param("id")
	if err := h.bookService.Delete(context.Background(), id); err != nil {
		h.logger.Error().Err(err).Str("book_id", id).Msg("Failed to delete book")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	h.logger.Info().Str("book_id", id).Msg("Book deleted successfully")
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) BorrowBook(c echo.Context) error {
	claimsInterface := c.Get("user")
	claims, ok := claimsInterface.(jwt.MapClaims)
	if !ok {
		h.logger.Error().Msg("Cannot parse user claims from token")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "cannot parse user claims"})
	}

	// Берём UUID пользователя из "sub"
	userIDStr, ok := claims["sub"].(string)
	if !ok {
		h.logger.Error().Msg("Cannot parse user ID from token")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "cannot parse user ID"})
	}

	// Получаем книгу по ID
	id := c.Param("id")
	book, err := h.bookService.GetByID(context.Background(), id)
	if err != nil {
		h.logger.Warn().Str("book_id", id).Msg("Book not found for borrow")
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	if !book.Available {
		h.logger.Warn().Str("book_id", id).Msg("Book already borrowed")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "already borrowed"})
	}

	// Помечаем книгу как занятое
	book.Available = false
	if err := h.bookService.Update(context.Background(), book); err != nil {
		h.logger.Error().Err(err).Str("book_id", id).Msg("Failed to borrow book")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Создаём запись в borrows
	borrow := &models.Borrow{
		UserID:     uuid.MustParse(userIDStr),
		BookID:     book.ID,
		BorrowedAt: time.Now(),
	}

	if err := h.bookService.CreateBorrow(context.Background(), borrow); err != nil {
		h.logger.Error().Err(err).Msg("Failed to create borrow record")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	h.logger.Info().Str("book_id", id).Str("user_id", userIDStr).Msg("Book borrowed successfully")
	return c.JSON(http.StatusOK, map[string]string{"status": "borrowed"})
}
