package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/user"
	"strconv"
	"time"

	"github.com/aquasecurity/table"
)

type Todo struct{
	User string
	Title string
	Completed bool
	CreatedAt time.Time
	CompletedAt *time.Time
	IPAddress string
}

type Todos []Todo

func (todos *Todos) validateIndex(index int) error{
	// Validation Task
	if index < 0 || index >= len(*todos) {
		err := errors.New("Invalid index")
		fmt.Println(err)
		return err
	}
	
	return nil
}


func (todos *Todos) add(title string) {
    // Add Task
    currentUser, err := user.Current()
    if err != nil {
        fmt.Println("Error fetching the user:", err)
        return
    }

    var ipAddress string
    ifaces, err := net.Interfaces()
    if err != nil {
        fmt.Println("Error getting network interfaces:", err)
        return
    }

    for _, i := range ifaces {
        addrs, err := i.Addrs()
        if err != nil {
            fmt.Println("Error getting addresses:", err)
            continue
        }
        for _, addr := range addrs {
            var ip net.IP
            switch v := addr.(type) {
            case *net.IPNet:
                ip = v.IP
            case *net.IPAddr:
                ip = v.IP
            }
            // Skip loopback and IPv6 addresses
            if ip.IsLoopback() || ip.To4() == nil {
                continue
            }
            ipAddress = ip.String() // Save the first valid IPv4 address
            break
        }
        if ipAddress != "" {
            break
        }
    }

    todo := Todo{
        User:        currentUser.Username,
        Title:       title,
        Completed:   false,
        CompletedAt: nil,
        CreatedAt:   time.Now(),
        IPAddress:   ipAddress, // Assign the IP address here
    }

    *todos = append(*todos, todo)
}


func (todos *Todos) delete(index int) error{
	// Delete Task
	t := *todos

	if err := t.validateIndex(index); err != nil{
		return err
	}

	*todos = append(t[:index], t[index+1:]...)
	return nil
}

func (todos * Todos) toggle(index int) error{
	// Done the task
	t := *todos

	if err := t.validateIndex(index); err != nil {
		return err
	}

	isCompleted := t[index].Completed

	if !isCompleted {
		completionTime := time.Now()
		t[index].CompletedAt = &completionTime
	}

	t[index].Completed = !isCompleted

	return nil
}



func (todos * Todos) edit(index int, title string) error{
	// Edit the task
	t := *todos

	if err := t.validateIndex(index); err != nil {
		return err
	}

	t[index].Title = title

	return nil
}

func (todos *Todos) print() {
    table := table.New(os.Stdout)
    table.SetRowLines(false)
    table.SetHeaders("#", "Title", "Completed", "Created At", "Completed At", "User", "IP Address") // Add IP Address header
    for index, t := range *todos {
        completed := "❌"
        completedAt := ""

        if t.Completed {
            completed = "✔️"
            if t.CompletedAt != nil {
                completedAt = t.CompletedAt.Format(time.RFC1123)
            }
        }

        table.AddRow(
            strconv.Itoa(index),
            t.Title,
            completed,
            t.CreatedAt.Format(time.RFC1123),
            completedAt,
            t.User,
            t.IPAddress, // Add IP Address here
        )
    }
    table.Render()
}
