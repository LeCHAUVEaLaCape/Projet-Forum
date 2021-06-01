// var for buttons
var modify_photo = document.getElementsByClassName("modify")[0]
var modify_age = document.getElementsByClassName("modify")[1]
var modify_address = document.getElementsByClassName("modify")[2]
var modify_fewWords = document.getElementsByClassName("modify")[3]



// var for photo
var photo_modifier = document.getElementsByClassName("photo-container")[0]
modify_photo.addEventListener("click", () => {
    photo_modifier.classList.toggle("photo-modifier-visible")
})

// var for age
var age_info = document.getElementsByClassName("age-database")[0]
var age_modifier = document.getElementsByClassName("input-age")[0]
modify_age.addEventListener("click", () => {
    age_info.classList.toggle("age-hidden")
    age_modifier.classList.toggle("age-modifier-visible")
})

// var for address
var address_info = document.getElementsByClassName("address-database")[0]
var address_modifier = document.getElementsByClassName("input-address")[0]
modify_address.addEventListener("click", () => {
    address_info.classList.toggle("address-hidden")
    address_modifier.classList.toggle("address-modifier-visible")
})

// var for fewWords
var fewWords_info = document.getElementsByClassName("fewWords-database")[0]
var fewWords_modifier = document.getElementsByClassName("input-fewWords")[0]
modify_fewWords.addEventListener("click", () => {
    fewWords_info.classList.toggle("fewWords-hidden")
    fewWords_modifier.classList.toggle("fewWords-modifier-visible")
})


// Affichage des boutons pour supprimer le compte
var show_delete = document.getElementById("show-delete-account")
var hide_delete = document.getElementById("hide-delete-account")
var delete_account_container = document.getElementsByClassName("del-account-container")[0]

show_delete.addEventListener("click", () => {
    delete_account_container.classList.add("del-active")
})
hide_delete.addEventListener("click", () => {
    delete_account_container.classList.remove("del-active")
})
