ws = null;

function deleteItem(id) {
  if (ws !== null) {
    ws.send(JSON.stringify({ action: "delete", data: id }));
  }
}

function completeItem(id) {
  if (ws !== null) {
    ws.send(JSON.stringify({ action: "complete", data: id }));
  }
}

function DisplayNewTodoItem(id, title, completed) {
  list.innerHTML +=
    `<tr id="${id}">
      <td>
        <kbd>
          <strong${completed ? ' style="text-decoration: line-through;"' : ""}>
            ${title}
          </strong>
        </kbd>
      </td>
      <td>
        <input type="checkbox" role="switch" onclick="completeItem('${id}')"${completed ? " disabled checked" : ""}/>
      </td>
      <td><button onclick="deleteItem('${id}')">Delete</button></td>
    </tr>`;
}

document.addEventListener('DOMContentLoaded', (event) => {
  list = document.getElementById("todo-list");
  ws = new WebSocket(
    location.protocol === "https" ? "wss://" : "ws://" +
      location.host +
      "/ws/todo"
  );

  ws.onopen = function () {
    ws.send('{ "action": "get_list" }');
  };

  ws.onmessage = function (event) {
    d = JSON.parse(event.data);

    switch (d.answer) {
      case "todo_list":
        if (d.data !== null) {
          for (let i = 0; i < d.data.length; i++) {
            DisplayNewTodoItem(d.data[i].id, d.data[i].title, d.data[i].completed);
          }
        }
        break;
      case "added_item":
        DisplayNewTodoItem(d.data.id, d.data.title, d.data.completed);
        break;
      case "completed_item":
        document.getElementById(d.data).children[0].children[0].children[0].setAttribute("style", "text-decoration: line-through;");
        document.getElementById(d.data).children[1].children[0].setAttribute("checked", "true");
        document.getElementById(d.data).children[1].children[0].setAttribute("disabled", "true");
        break;
      case "deleted_item":
        document.getElementById(d.data).remove();
        break;
      default:
        console.log("WebSocket message received:", event.data);
    }
  };

  document.getElementById("create_button").addEventListener("click", function (event) {
    ws.send(
      JSON.stringify({
        action: "create",
        data: { title: document.getElementById("new_title").value },
      })
    );
    document.getElementById("new_title").value = "";
  });

  document.getElementById("new_title").addEventListener("keypress", function (event) {
    if (event.key === "Enter") {
      event.preventDefault();
      document.getElementById("create_button").click();
    }
  });
});
