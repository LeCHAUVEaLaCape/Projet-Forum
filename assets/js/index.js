var order =  document.getElementsByClassName("all-posts")[0]
var change_text_post = document.getElementById("post-text")

var newest = document.getElementById("newest")
var latest = document.getElementById("latest")

newest.addEventListener("click", () => {
    order.style.flexWrap = "wrap"
    change_text_post.innerHTML = "Nouveaux Posts"
})
latest.addEventListener("click", () => {
    order.style.flexWrap = "wrap-reverse"
    change_text_post.innerHTML = "Anciens Posts"
})

var grand = document.getElementById("grand")
var moyen = document.getElementById("moyen")
var compact = document.getElementById("compact")

var message = document.getElementsByClassName("message")
var posts_container = document.getElementsByClassName("singlePost-container")

grand.addEventListener("click", () => {
    for (i = 0 ; i < posts_container.length ; i++) {
        posts_container[i].style.width = "70%"
        posts_container[i].style.margin = "10px 0px 25px 0px"
        message[i].style.height = "175px"
    }
})

moyen.addEventListener("click", () => {
    for (i = 0 ; i < posts_container.length ; i++) {
        posts_container[i].style.width = "38%"
        posts_container[i].style.margin = "10px 20px 25px 20px"
        message[i].style.height = "130px"
    }
})