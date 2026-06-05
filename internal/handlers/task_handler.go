package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"task-manager-api/internal/models"
	"task-manager-api/internal/services"

	"github.com/gin-gonic/gin"
)

// TaskHandler обрабатывает запросы задач.
type TaskHandler struct {
	taskService *services.TaskService
}

// NewTaskHandler создает обработчик задач.
func NewTaskHandler(taskService *services.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

// CreateTask создает задачу.
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req models.CreateTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID := c.MustGet("userID").(uint)

	task, err := h.taskService.CreateTask(
		c.Request.Context(),
		userID,
		req,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// ListTasks возвращает список задач.
func (h *TaskHandler) ListTasks(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}
	status := c.Query("status")

	response, err := h.taskService.ListTasks(
		c.Request.Context(),
		userID,
		page,
		limit,
		status,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetTask получает задачу.
func (h *TaskHandler) GetTask(c *gin.Context) {
	taskID, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "некорректный идентификатор",
		})
		return
	}

	userID := c.MustGet("userID").(uint)

	task, err := h.taskService.GetTask(
		c.Request.Context(),
		uint(taskID),
		userID,
	)

	if err != nil {
		if errors.Is(err, services.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "forbidden",
			})
			return
		}

		if h.taskService.IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "задача не найдена",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, task)
}

// UpdateTask обновляет задачу.
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	taskID, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "некорректный идентификатор",
		})
		return
	}

	var req models.UpdateTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID := c.MustGet("userID").(uint)

	task, err := h.taskService.UpdateTask(
		c.Request.Context(),
		uint(taskID),
		userID,
		req,
	)

	if err != nil {
		if errors.Is(err, services.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "forbidden",
			})
			return
		}

		if h.taskService.IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "задача не найдена",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, task)
}

// DeleteTask удаляет задачу.
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	taskID, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "некорректный идентификатор",
		})
		return
	}

	userID := c.MustGet("userID").(uint)

	err = h.taskService.DeleteTask(
		c.Request.Context(),
		uint(taskID),
		userID,
	)

	if err != nil {
		if errors.Is(err, services.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "forbidden",
			})
			return
		}

		if h.taskService.IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "задача не найдена",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// UpdateTaskStatus обновляет статус задачи.
func (h *TaskHandler) UpdateTaskStatus(c *gin.Context) {
	taskID, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "некорректный идентификатор",
		})
		return
	}

	var req models.UpdateTaskStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID := c.MustGet("userID").(uint)

	task, err := h.taskService.UpdateStatus(
		c.Request.Context(),
		uint(taskID),
		userID,
		req.Status,
	)

	if err != nil {
		if errors.Is(err, services.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "forbidden",
			})
			return
		}

		if h.taskService.IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "задача не найдена",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, task)
}
