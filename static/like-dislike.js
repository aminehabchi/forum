document.addEventListener("click", (e) => {
  const button = e.target.closest(".likedislikebtn");
  if (button) {
    likeDislike(button.dataset.id, button.dataset.action, button.dataset.type);
  }
});

const likeDislike = (Id, action, type) => {
  fetch(`/like-dislike?action=${action}&commentid=${Id}&type=${type}`)
    .then((res) => {

      if (res.status == 200) {
        const likespan = document.querySelector(`#like_post-${Id}`);
        const dislikespan = document.querySelector(`#dislike_post-${Id}`);
        const likebtn = likespan.parentElement;
        const dislikebtn = dislikespan.parentElement;

        if (action == "like") {
          
          if (!likebtn.classList.contains("liked-btn")) {
            likespan.textContent = Number(likespan.textContent) + 1;
            likebtn.classList.add("liked-btn");

            if (dislikebtn.classList.contains("liked-btn")) {
              dislikespan.textContent = Number(dislikespan.textContent) - 1;
              dislikebtn.classList.remove("liked-btn");
            }
          } else {
            likespan.textContent = Number(likespan.textContent) - 1;
            likebtn.classList.remove("liked-btn");
          }
        } else {

          if (!dislikebtn.classList.contains("liked-btn")) {
            dislikespan.textContent = Number(dislikespan.textContent) + 1;
            dislikebtn.classList.add("liked-btn");

            if (likebtn.classList.contains("liked-btn")) {
              likespan.textContent = Number(likespan.textContent) - 1;
              likebtn.classList.remove("liked-btn");
            }
          } else {
            dislikespan.textContent = Number(dislikespan.textContent) - 1;
            dislikebtn.classList.remove("liked-btn");
          }
        }
      } else if (res.status == 403) {
        window.location.replace("/login");
      } else {
        window.location.replace(`/error?s=${res.status}&st=${res.statusText}`);
      }
    })
    .catch((err) => {
      console.log(err);
    });
};
