var notif_btn = document.getElementById("notif-btn")
var notif_container = document.getElementsByClassName("notif-container")[0]
var notif = document.getElementsByClassName("notif")
var all_notif = document.getElementById("all-notif")

if (localStorage.getItem("notif_seen") == null) {
    localStorage.setItem("notif_seen", 0)
}

if (all_notif.childElementCount > localStorage.getItem("notif_seen")){
    notif_btn.setAttribute("src", "../assets/images/notification-active.png")
} else {
    notif_btn.setAttribute("src", "../assets/images/notification-inactive.png")
}

notif_btn.addEventListener("click", () => {
    notif_container.classList.toggle("notif-active")
    notif_btn.setAttribute("src", "../assets/images/notification-inactive.png")
    localStorage.setItem("notif_seen", all_notif.childElementCount)
})

var action = document.getElementsByClassName("action")
for (let i = 0 ; i < action.length ; i++) {
    if (action[i].innerHTML == "like") {
        action[i].innerHTML = "aimé"
    } else if (action[i].innerHTML == "comment"){
        action[i].innerHTML = "commenté"
    }
}