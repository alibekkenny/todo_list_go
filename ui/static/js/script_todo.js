colorIcon(0);

tooltip();
let id = 6;
let idCard;
let idArray;
const myModal = new bootstrap.Modal(document.getElementById("modal-card"));

let arrayCard = [
    ["cel1", "0"],
    ["cel2", "0"],
    ["cel3", "1"],
    ["cel4", "1"],
    ["cel5", "1"]
];
let arrayType = [
    ["dev", "#2d86ba", "#ceecfd"],
    ["uni", "#13854e", "#d6ede2"]
];

chargeType();

function changeDemoText(ev) {
    document.getElementById("demoType").innerText = ev.target.value;
}

var TxtColorDemo = document.getElementById("color-for-texte");
TxtColorDemo.addEventListener(
    "input",
    function() {
        document.getElementById("demoType").style.color = TxtColorDemo.value;
    },
    false
);

var bckGroundDemo = document.getElementById("bck-color");
bckGroundDemo.addEventListener(
    "input",
    function() {
        document.getElementById("demoType").style.backgroundColor =
            bckGroundDemo.value;
    },
    false
);

function addType() {
    name = document.getElementById("nameType").value;
    bckColor = document.getElementById("bck-color").value;
    TxtColor = document.getElementById("color-for-texte").value;
    if (name.trim() != "" && bckColor.trim() != "" && TxtColor.trim() != "") {
        newType = [name, TxtColor, bckColor];
        arrayType.push(newType);
        document.getElementById("nameType").value = "";
        document.getElementById("demoType").innerText = "dev";
    } else {
        let toast = document.getElementById("toast");
        toast.innerHTML += `
    <div class="position-fixed bottom-0 end-0 p-3" style="z-index: 11">
      <div id="liveToast" class="toast" role="alert" aria-live="assertive" aria-atomic="true">
        <div class="toast-header">
          <strong class="me-auto">Nouveau type</strong>
          <button type="button" class="btn-close" data-bs-dismiss="toast" aria-label="Close"></button>
        </div>
        <div class="toast-body">
          Les champs ne sont pas bien remplis !
        </div>
      </div>
    </div>
    `;
        var toastLiveExample = document.getElementById("liveToast");
        var toastBtstrap = new bootstrap.Toast(toastLiveExample);
        toastBtstrap.show();
    }
}

function colorIcon(id) {
    let leftNav = document.getElementsByClassName("left-nav-item")[0];
    icon = leftNav.getElementsByClassName("left-nav-icon");
    for (let element of icon) {
        element.classList.remove("active-icon-menu");
    }
    icon[id].classList.add("active-icon-menu");
}

function loadPage(id) {
    page = document.getElementsByClassName("main");
    for (let element of page) {
        element.classList.add("d-none");
    }
    page[id].classList.remove("d-none");
    colorIcon(id);
}

// function chargeType() {
//   arrayCard.forEach((element) => {
//     let card = document.getElementById(element[0]);
//     console.log(card.getElementsByClassName("badge"));
//     if (card.getElementsByClassName("badge")[0] != undefined) {
//       card.getElementsByClassName("badge")[0].remove();
//     }
//     let children = card.children;
//     if (element[1] != -1) {
//       type = arrayType[element[1]];
//       console.log(type);
//       let colorBck = type[2];
//       let colorTxt = type[1];
//       let txt = type[0];
//       children[0].innerHTML += `<span class="badge" style='background-color:${colorBck};color:${colorTxt};'>${txt}</span>`;
//     }
//   });
// }

function getIdArrayCard(id) {
    let cpt = 0;
    save = null;
    arrayCard.forEach((element) => {
        if (element[0] == id) {
            save = cpt;
        }
        cpt++;
    });
    return save;
}

function createOption(id) {
    console.log(arrayCard);
    let cpt = 0;
    sel = document.getElementById("select-type");
    sel.innerHTML = "";
    cardIdType = arrayCard[id][1];
    if (cardIdType == -1) {
        sel.innerHTML += `<option value=-1 selected>Choisir un type</option>`;
    } else {
        sel.innerHTML += `<option value=-1 selected>smth</option>`;
    }

    arrayType.forEach(async(element) => {
        if (cardIdType == cpt) {
            sel.innerHTML += `<option value=${cpt} selected>${element[0]}</option>`;
        } else {
            sel.innerHTML += `<option value=${cpt}>${element[0]}</option>`;
        }

        cpt++;
    });
}

async function openModalCard(ev) {
    title = ev.currentTarget.getElementsByTagName("h6")[0].innerText;
    text = ev.currentTarget.getElementsByTagName("span")[0].innerText;

    document.getElementById("modal-title").innerText = title;
    document.getElementById("modal-content").innerText = text;
    idCard = ev.currentTarget.id;
    idArray = await getIdArrayCard(idCard);
    createOption(idArray);
    myModal.show();
}

function saveModifCard() {
    title = document.getElementById("modal-title").innerText;
    text = document.getElementById("modal-content").innerText;

    card = document.getElementById(idCard);
    card.getElementsByTagName("h6")[0].innerText = title;
    card.getElementsByTagName("span")[0].innerText = text;

    smallText = text.substr(0, 50);
    card.getElementsByClassName("small-content-text")[0].innerText =
        smallText + "...";
    myModal.hide();

    arrayCard[idArray][1] = document.getElementById("select-type").value;
    console.log(idArray);
    console.log(arrayCard);
    chargeType();
}

function addCard(ev) {
    ev.currentTarget.outerHTML =
        `<div class='card card-col' draggable='true' ondragstart='drag(event)' id='cel${id}' onclick='openModalCard(event)'>
    <div class='card-body'>
    <div class="card-action">
                    <h6 class="card-title" contenteditable onclick="event.stopPropagation();">
                      Card title
                    </h6>
                    <div class="dropdown">
                      <i class="fas fa-ellipsis-v" id="dropdownMenuButton1" data-bs-toggle="dropdown" aria-expanded="false" onclick="event.stopPropagation();"></i>
                      <ul class="dropdown-menu" aria-labelledby="dropdownMenuButton1">
                        <li><div class="delet-btn" onclick="deletCard(event)">&nbsp;<i class="fas fa-trash-alt i"></i> Delete</div></li>
                      </ul>
                    </div>
                  </div>
                  <span class="card-text fs-6 full-content-text">Some quick example text to build on the card title and make up the bulk of the card's content.</span>
    <p class='small-content-text'>Some quick example text to build on the card title...</p>
  </div>
</div>` + ev.currentTarget.outerHTML;
    newCard = ["cel" + id, "-1"];
    id++;
    tooltip();

    arrayCard.push(newCard);
}

function deletCard(ev) {
    ev.stopPropagation();
    element = ev.target.parentNode;
    element1 = element.parentNode;
    element2 = element1.parentNode;
    element3 = element2.parentNode;
    cardbody = element3.parentNode;
    card = cardbody.parentNode;
    console.log(card);
    card.remove();
}

function allowDrop(ev) {
    ev.preventDefault();
}

function drag(ev) {
    ev.dataTransfer.setData("Text", ev.target.id);
}

function drop(ev) {
    var data = ev.dataTransfer.getData("Text");
    let elem = ev.target;
    if (elem.tagName != "DIV" || elem.id == "") {
        while (elem.tagName != "DIV" || elem.id == "") {
            elem = elem.parentNode;
        }
    }
    if (elem.tagName == "DIV" && elem.id != "") {
        if (elem.classList.contains("col")) {
            elem.insertBefore(
                document.getElementById(data),
                elem.childNodes[elem.childNodes.length - 2]
            );
        } else {
            elem.after(document.getElementById(data));
        }
        ev.preventDefault();
    } else {
        console.log("mauvais positionnement");
    }
}

function tooltip() {
    var popoverTriggerList = [].slice.call(
        document.querySelectorAll('[data-bs-toggle="popover"]')
    );
    var popoverList = popoverTriggerList.map(function(popoverTriggerEl) {
        return new bootstrap.Popover(popoverTriggerEl);
    });
}