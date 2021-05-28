var order = document.getElementsByClassName("all-posts")[0];
var change_text_post = document.getElementById("post-text");

var newest = document.getElementById("newest");
var latest = document.getElementById("latest");

newest.addEventListener("click", () => {
  order.style.flexWrap = "wrap";
  change_text_post.innerHTML = "Anciens Posts";
});
latest.addEventListener("click", () => {
  order.style.flexWrap = "wrap-reverse";
  change_text_post.innerHTML = "Nouveaux Posts";
});

var grand = document.getElementById("grand");
var moyen = document.getElementById("moyen");
var compact = document.getElementById("compact");

var message = document.getElementsByClassName("message");
var posts_container = document.getElementsByClassName("singlePost-container");
var categories = document.getElementsByClassName("categories");
var each_cat = document.getElementsByClassName("each-cat");

var body_post = document.getElementsByClassName("body-post");
var see_more = document.getElementsByClassName("see-more")

var time = document.getElementsByClassName("time-posted")
var author = document.getElementsByClassName("author-time")
var color = document.getElementsByClassName("color")

grand.addEventListener("click", () => {
  for (let i = 0; i < posts_container.length; i++) {
    posts_container[i].style.width = "70%";
    posts_container[i].style.margin = "10px 0px 25px 0px";
    message[i].style.height = "175px";
    message[i].style.borderRadius = "10px 10px 0px 0px";
    see_more[i].style.borderRadius = "0px 0px 10px 10px";
    author[i].style.display = "initial"
  }
  for (let i = 0; i < categories.length; i++) {
    categories[i].style.flexWrap = "wrap";
    body_post[i].style.flexWrap = "wrap"
    see_more[i].style.width = "100%"
  }
  for (let i = 0; i < each_cat.length; i++) {
    color[i].style.height = "27px";
    each_cat[i].style.display = "initial";
  }
});

moyen.addEventListener("click", () => {
  for (i = 0; i < posts_container.length; i++) {
    posts_container[i].style.width = "38%";
    posts_container[i].style.margin = "10px 20px 25px 20px";
    message[i].style.height = "130px";
    message[i].style.borderRadius = "10px 10px 0px 0px";
    see_more[i].style.borderRadius = "0px 0px 10px 10px";
    author[i].style.display = "initial"
  }
  for (let i = 0; i < categories.length; i++) {
    categories[i].style.flexWrap = "wrap";
    body_post[i].style.flexWrap = "wrap"
    see_more[i].style.width = "100%"
  }
  for (let i = 0; i < each_cat.length; i++) {
    color[i].style.height = "27px";
    each_cat[i].style.display = "initial";
  }
});

compact.addEventListener("click", () => {
  for (i = 0; i < posts_container.length; i++) {
    posts_container[i].style.width = "70%";
    posts_container[i].style.margin = "0px";
    message[i].style.height = "32px";
    message[i].style.borderRadius = "0px";
    see_more[i].style.borderRadius = "0px";
    author[i].style.display = "flex"
  }
  for (let i = 0; i < categories.length; i++) {
    categories[i].style.flexWrap = "initial";
    body_post[i].style.flexWrap = "initial"
    see_more[i].style.width = "9%"
  }
    for (let i = 0; i < each_cat.length; i++) {
      color[i].style.height = "16px";
      each_cat[i].style.display = "none";
  }
});
