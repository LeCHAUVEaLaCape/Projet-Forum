
// Modif the MAIN post
var display_modif = document.getElementById("display-modif")
var post = document.getElementsByClassName("message")[0]
var modif = document.getElementsByClassName("form-modif")[0]
var textarea_modif = document.getElementById("modif")

// make the modif visible/invisible
if (display_modif != null) {
    display_modif.addEventListener("click", () => {
        post.classList.toggle("hide-message")
        modif.classList.toggle("see-modif")
    })

    // make the text is the textarea correct
    new_text = textarea_modif.innerHTML.replace(/&lt;br&gt;/g, '').replace(/<br>/g, '')
    textarea_modif.innerHTML = new_text
}


// Modif the comment
var display_modif_comment = document.getElementsByClassName("display-modif-comment")
var comment = document.getElementsByClassName("comment")
var modif_comment = document.getElementsByClassName("form-modif-comment")
var textarea_modif_comment = document.getElementsByClassName("modifComment")
var resize_comment_box = document.getElementsByClassName("form-modif-comment")

var comment_container = document.getElementsByClassName("comment-container")
            

// change the number of posts
console.log(comment_container.length)
document.getElementById("nb-comment").innerHTML = comment_container.length

console.log(display_modif_comment[0].parentNode.children[0])
for (let i = 0 ; i < display_modif_comment.length ; i++) {
    display_modif_comment[i].addEventListener("click", () => {
        display_modif_comment[i].parentNode.children[0].classList.toggle("hide-message")
        modif_comment[i].classList.toggle("see-modif")
        resize_comment_box[i].classList.toggle("resize-box")
    })
    // make the text is the textarea correct
    new_text = textarea_modif_comment[i].innerHTML.replace(/&lt;br&gt;/g, '').replace(/<br>/g, '')
    textarea_modif_comment[i].innerHTML = new_text
}
