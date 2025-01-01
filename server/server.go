package server

import (
	models "CRUD/model"
	"CRUD/model/db"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type server struct {
	mutex       sync.Mutex
	taskStorage db.TaskStorage
}

func NewServer(taskStorage db.TaskStorage) *server {
	return &server{sync.Mutex{}, taskStorage}
}

func (server *server) getAllTasksHandler(respWriter http.ResponseWriter, request *http.Request) {
	server.mutex.Lock()
	tasks, err := server.taskStorage.GetAllTasks()
	defer server.mutex.Unlock()
	if err != nil {
		sendErr := sendError(&respWriter, err, http.StatusInternalServerError)
		if sendErr != nil {
			fmt.Println(sendErr.Error())
		}
	}
	body, serializeErr := json.Marshal(tasks)
	if serializeErr != nil {
		sendErr := sendError(&respWriter, serializeErr, http.StatusInternalServerError)
		if sendErr != nil {
			fmt.Println(sendErr.Error())
		}
	}
	sendErr := sendResult(&respWriter, body)
	if sendErr != nil {
		fmt.Println(sendErr.Error())
	}
}

func (server *server) getTaskByIdHandler(respWriter http.ResponseWriter, request *http.Request) {
	taskId, err := strconv.ParseUint(mux.Vars(request)["id"], 10, 64)
	if err != nil {
		sendErr := sendError(&respWriter, err, http.StatusBadRequest)
		if sendErr != nil {
			fmt.Println(sendErr.Error())
		}
	}
	server.mutex.Lock()
	task, dbErr := server.taskStorage.GetTaskById(taskId)
	defer server.mutex.Unlock()
	if dbErr != nil {
		sendErr := sendError(&respWriter, dbErr, http.StatusInternalServerError)
		if sendErr != nil {
			fmt.Println(sendErr.Error())
		}
	}
	body, serializeErr := json.Marshal(task)
	if serializeErr != nil {
		sendErr := sendError(&respWriter, serializeErr, http.StatusInternalServerError)
		if sendErr != nil {
			fmt.Println(sendErr.Error())
		}
	}
	sendErr := sendResult(&respWriter, body)
	if sendErr != nil {
		fmt.Println(sendErr.Error())
	}
}

func (server *server) createTaskHandler(respWriter http.ResponseWriter, request *http.Request) {
	body, bodyErr := ioutil.ReadAll(request.Body)
	if bodyErr != nil {
		sendErr := sendError(&respWriter, bodyErr, http.StatusBadRequest)
		if sendErr != nil {
			fmt.Println(sendErr.Error())
		}
	}
	var task models.Task
	deserializeErr := json.Unmarshal(body, &task)
	if deserializeErr != nil {
		sendErr := sendError(&respWriter, deserializeErr, http.StatusBadRequest)
		if sendErr != nil {
			fmt.Println(sendErr.Error())
		}
	}
	server.mutex.Lock()
	taskId, dbErr := server.taskStorage.CreateTask(&task)
	defer server.mutex.Unlock()
	if dbErr != nil {
		sendErr := sendError(&respWriter, dbErr, http.StatusInternalServerError)
		if sendErr != nil {
			fmt.Println(sendErr.Error())
		}
	}
	serId := []byte(string(taskId))
	sendErr := sendResult(&respWriter, serId)
	if sendErr != nil {
		fmt.Println(sendErr.Error())
	}

}

func (server *server) updateTaskHandler(respWriter http.ResponseWriter, request *http.Request) {
	taskId, err := strconv.ParseUint(mux.Vars(request)["id"], 10, 64)
	if err != nil {
		sendErr := sendError(&respWriter, err, http.StatusBadRequest)
		if sendErr != nil {
			fmt.Println(sendErr.Error())
		}
	}
	body, bodyErr := ioutil.ReadAll(request.Body)
	if bodyErr != nil {
		sendErr := sendError(&respWriter, bodyErr, http.StatusBadRequest)
		if sendErr != nil {
			fmt.Println(sendErr.Error())
		}
	}
	var task models.Task
	deserializeErr := json.Unmarshal(body, &task)
	if deserializeErr != nil {
		sendErr := sendError(&respWriter, deserializeErr, http.StatusBadRequest)
		if sendErr != nil {
			fmt.Println(sendErr.Error())
		}
	}
	task.Id = taskId
	server.mutex.Lock()
	upCount, dbErr := server.taskStorage.UpdateTask(&task)
	defer server.mutex.Unlock()
	if dbErr != nil {
		sendErr := sendError(&respWriter, dbErr, http.StatusInternalServerError)
		if sendErr != nil {
			fmt.Println(sendErr.Error())
		}
	}
	serUpCount := []byte(string(upCount))
	sendErr := sendResult(&respWriter, serUpCount)
	if sendErr != nil {
		fmt.Println(sendErr.Error())
	}
}

func (server *server) deleteTaskHandler(respWriter http.ResponseWriter, request *http.Request) {
	taskId, err := strconv.ParseUint(mux.Vars(request)["id"], 10, 64)
	if err != nil {
		sendErr := sendError(&respWriter, err, http.StatusBadRequest)
		if sendErr != nil {
			fmt.Println(sendErr.Error())
		}
	}
	server.mutex.Lock()
	delCount, dbErr := server.taskStorage.DeleteTask(taskId)
	defer server.mutex.Unlock()
	if dbErr != nil {
		sendErr := sendError(&respWriter, dbErr, http.StatusInternalServerError)
		if sendErr != nil {
			fmt.Println(sendErr.Error())
		}
	}
	serDelCount := []byte(string(delCount))
	sendErr := sendResult(&respWriter, serDelCount)
	if sendErr != nil {
		fmt.Println(sendErr.Error())
	}
}

func (server *server) buildRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/tasks", server.createTaskHandler).Methods("POST")
	router.HandleFunc("/tasks", server.getAllTasksHandler).Methods("GET")
	router.HandleFunc("/tasks/{id:[0-9]+}", server.getTaskByIdHandler).Methods("GET")
	router.HandleFunc("/tasks/{id:[0-9]+}", server.updateTaskHandler).Methods("PUT")
	router.HandleFunc("/tasks/{id:[0-9]+}", server.deleteTaskHandler).Methods("DELETE")
	return router
}

func (server *server) Start() error {
	host := os.Getenv("SERVER_HOST")
	if host == "" {
		return errors.New("$SERVER_HOST not set")
	}
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		return errors.New("$SERVER_PORT not set")
	}
	router := server.buildRouter()
	if router == nil {
		return errors.New("Building Router failed")
	}
	err := http.ListenAndServe(host+":"+port, router)
	if err != nil {
		return errors.New("ListenAndServe: " + err.Error())
	}
	return nil
}
