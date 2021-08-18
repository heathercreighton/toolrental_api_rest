package tools

import(
  "encoding/json"
  "fmt"
  "sort"
  "strconv"
  "strings"
  "net/http"
  db "main/db"
)

type Tools struct{
  db *db.Tool
}

func Handlers() *Tools{
  return &Tools{db: &db.Tool{}}
}


func (t *Tools) Root(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type","application/json")
  
  switch r.Method {

  case "GET":  
    tools, err := t.db.All()
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }
   
    q, ok := r.URL.Query()["sort"]                      // 1
		if ok {                                             // 2
			order := q[0]                                     // 3
			if order == "asc" {                               // 4
				sort.SliceStable(tools, func(i, j int) bool {   // 5
					return tools[i].Price < tools[j].Price
				})

			} else if order == "desc" {
        sort.SliceStable(tools, func(i, j int) bool {   // 5
					return tools[i].Price > tools[j].Price
				})
      }
		}
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(tools)



  case"POST":
    tool := &db.Tool{}

    if err := json.NewDecoder(r.Body).Decode(tool); err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }

    if err := t.db.Create(tool); err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(tool)

  }

}

func (t *Tools) Items(w http.ResponseWriter, r *http.Request) {
	// We can do all the prework before checking for the HTTP Method
  
  idParam := strings.TrimPrefix(r.URL.Path, "/api/tools/") //trim the ID from the route
  id, err  := strconv.Atoi(idParam)   //converts the ID to an integer

  if err != nil {
      http.Error(w, err.Error(), http.StatusBadRequest)
      return
    }
  tool, err := t.db.FindByID(id) //finds the tool in the database
  if err != nil || tool == nil {
        http.Error(w, "Invalid ID. Please try again.", http.StatusBadRequest)
        return
  }

  switch r.Method{

    case "GET": 
      json.NewEncoder(w).Encode(tool)

    case "PUT", "PATCH":
      tool := &db.Tool{}
      if err := json.NewDecoder(r.Body).Decode(tool); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
      }

      if err := t.db.Update(id, tool); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
      }

      w.WriteHeader(http.StatusOK)
      json.NewEncoder(w).Encode(tool)

    case "DELETE":
      if err := t.db.Delete(id); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
      }
      w.WriteHeader(http.StatusOK)
  }
}

func (t *Tools) ThrowError(w http.ResponseWriter, r *http.Request) {
  err := fmt.Errorf("There's a problem")
  if err != nil {
		http.NotFound(w,r)
		return // Don't forget to return or the function will attempt to keep going.
	}


}
