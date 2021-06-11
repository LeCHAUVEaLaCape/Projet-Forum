var see_more = document.getElementsByClassName("see-more")
var message = document.getElementsByClassName("message")
var see_more_txt = document.getElementsByClassName("see-more-txt")

var click = false

for (let i = 0 ; i < see_more.length ; i++) {
    see_more[i].addEventListener("click", () => {
        message[i].classList.toggle("message-full")
        if (click) {
            click = false
            see_more_txt[i].innerHTML = "See more..."
            console.log(click)
            console.log(see_more_txt.innerHTML)
        } else {
            click = true
            see_more_txt[i].innerHTML = "See less..."
            console.log(click)
            console.log(see_more_txt.innerHTML)
        }
    })
}