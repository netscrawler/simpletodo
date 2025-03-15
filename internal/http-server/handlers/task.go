package handlers

import (
	"net/http"
	"simpletodo/internal/model"
	"simpletodo/internal/service"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	taskService *service.TaskService
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

func (th *TaskHandler) Home(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", gin.H{})
}

func (th *TaskHandler) FetchTask(c *gin.Context) {
	t, err := th.taskService.GetAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "todoList", t)
}

func (th *TaskHandler) GetTaskForm(c *gin.Context) {
	c.HTML(http.StatusOK, "addTaskForm", gin.H{})
}

func (th *TaskHandler) AddTask(c *gin.Context) {
	title := c.PostForm("task")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task cannot be empty"})
		return
	}
	err := th.taskService.Create(c, title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	th.FetchTask(c)
}

func (th *TaskHandler) GetTaskUpdateForm(c *gin.Context) {
	id := c.Param("id")
	idx, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	task, err := th.taskService.GetByID(c, idx)
	if err != nil {
		if err == model.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "updateTaskForm", task)
}

func (th *TaskHandler) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	idx, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	askItem := c.PostForm("task")
	taskStatus := strings.ToLower(c.PostForm("done")) == "yes" || c.PostForm("done") == "on"

	err = th.taskService.Update(c, &model.Task{
		ID:    idx,
		Title: askItem,
		Done:  taskStatus,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	th.FetchTask(c)

}

func (th *TaskHandler) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	idx, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	err = th.taskService.Delete(c, idx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	th.FetchTask(c)
}
