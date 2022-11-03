// Show tasks in a list format
function showTasks(param) {
    fetch("api/tasks", {
        method: "GET",
        headers: {
            "Content-Type": "application/json"
        },
        body: null
    })
        .then(response => response.json())
        .then(json => {
            console.log(json)
            document.getElementById("myUl").innerHTML = "";
            json.data.forEach(function (item) {
                // var li = document.createElement("li");
                // var li = createCard(item);
                var li = createListItem(item);
                // var text = document.createTextNode("Task name: " + item.TaskName + ", Assignee: " + item.Assignee + ", Is done: " + "false");
                li.setAttribute('id', item.ID)
                // li.setAttribute('onclick','deleteTask(this);');
                // li.appendChild(text);
                document.getElementById('myUl').appendChild(li);
            });
        });
}

function createListItem(param) {
    let li = document.createElement('li');

    let div = document.createElement('div');
    div.className = 'form-check';

    label = document.createElement('label');
    label.className = 'form-check-label';
    label.type = 'text';
    label.textContent = param.Assignee == "" ? param.TaskName : param.TaskName + " Assigned to " + param.Assignee;
    if (param.isDone) label.style = "text-decoration: line-through;";

    input = document.createElement('input');
    input.className = 'checkbox';
    input.type = 'checkbox';
    input.checked = param.isDone;
    input.onclick = function () {
        toggleIsDone(param);
    };

    let i = document.createElement('i');
    i.className = 'input-helper';

    label.appendChild(input);
    label.appendChild(i);

    div.appendChild(label);

    let i2 = document.createElement('i');
    i2.className = 'material-icons';
    i2.textContent = 'close';
    i2.style = "position:absolute; right: 50px;";
    i2.onclick = function () {
        deleteTask(param);
    };

    let i3 = document.createElement('i');
    i3.className = 'material-icons';
    i3.textContent = 'edit';
    i3.style = "position:absolute; right: 85px;";
    i3.onclick = function () {
        console.log(param);
        editTask(param);
    };

    li.appendChild(div);
    li.appendChild(i2);
    li.appendChild(i3);
    return li;
}

function toggleIsDone(param) {
    param.isDone = !param.isDone;
    updateTask(param);
}

function updateTask(param) {
    console.log(param);
    fetch('api/tasks/' + param.id, {
        method: 'PATCH',
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ ID: param.ID, TaskName: param.TaskName, Assignee: param.Assignee, IsDone: param.isDone })
    }).then(response => {
        console.log(response);
        showTasks();
    }).catch(e => {
        console.log(e);
    });
}


// <div class="add-items d-flex">
//                            <input type="text" class="form-control todo-list-input" id="taskName" placeholder="e.g., Review design docs">
//                            <input type="text" class="form-control todo-list-input" id="assignee" placeholder="Assignee (optional)">
//                            <button onclick="createTask()" class="add btn btn-primary font-weight-bold todo-list-add-btn">Add</button> 
//                         </div>
function editTask(param) {
    let div = document.createElement('div');
    div.class = "add-items d-flex";

    let input = document.createElement('input');
    input.type = "text";
    input.className = "form-control todo-list-input";
    input.id = "editTaskName";
    input.placeholder = "e.g., Review design docs";
    input.value = param.TaskName;

    let input2 = document.createElement('input');
    input2.type = "text";
    input2.className = "form-control todo-list-input";
    input2.id = "editAssignee";
    input2.placeholder = "Assignee (optional)";
    input2.value = param.Assignee;

    let button = document.createElement('button');
    button.className = "add btn btn-primary font-weight-bold todo-list-add-btn";
    button.onclick = function () {
        updateTask({ ID: param.ID, TaskName: input.value, Assignee: input2.value, IsDone: param.isDone });
    }
    button.textContent = "Save";

    let button2 = document.createElement('button');
    button2.className = "add btn btn-secondary font-weight-bold todo-list-add-btn";
    button2.onclick = function () {
        showTasks();
    }
    button2.textContent = "Cancel";


    div.appendChild(input);
    div.appendChild(input2);
    div.appendChild(button);
    div.appendChild(button2);
    // button.onclick="editTask()"
    const listItem = document.getElementById(param.ID);
    listItem.parentNode.replaceChild(div, listItem);
}

function deleteTask(param) {
    fetch('api/tasks/' + param.ID, {
        method: 'DELETE',
        headers: {
            "Content-Type": "application/json"
        },
        body: null
    }).then(response => {
        console.log(response);
        showTasks();
    }).catch(e => {
        console.log(e);
    });
}

// Creates task in database 
function createTask() {
    console.log('here');
    var TaskName = document.getElementById('taskName').value;
    var Assignee = document.getElementById('assignee').value;
    fetch('api/tasks/', {
        method: 'POST',
        headers: {
            'Accept': 'application/json, text/plain, */*',
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ TaskName: TaskName, Assignee: Assignee })
    }).then(response => {
        console.log(response);
        showTasks();
        document.getElementById('taskName').value = "";
        document.getElementById('assignee').value = "";
    }).catch(e => {
        console.log(e);
    });
}

document.getElementById("add-task").addEventListener("click", createTask)