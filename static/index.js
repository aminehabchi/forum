document.querySelectorAll('.likedislikebtn').forEach(button => {
    button.addEventListener('click', () => {
        const id = button.dataset.id;
        const action = button.dataset.action;
        const type = button.dataset.type

        likeDislike(id, action, type);
    });
});

const likeDislike = (CommentId, action, type) => {

    fetch(`/like-dislike?action=${action}&commentid=${CommentId}&type=${type}`)
        .then(res => {

            if (!res.ok) {
                window.location.replace(`/error?s=${res.status}&st=${res.statusText}`)
                return
            }
            if (res.url.includes('/login')) {                
                window.location.replace('/login');
                return;
            }

            return res.json()

        })
        .then((data) => {

            document.querySelector(`#like_post-${CommentId}`).textContent = data.like
            document.querySelector(`#dislike_post-${CommentId}`).textContent = data.dislike
            document.querySelectorAll(`.likedislikebtn[data-id="${CommentId}"]`).forEach(button => {
                button.classList.remove('liked-btn');
                if (data.interaction === 1) {
                    document.querySelector(`.likedislikebtn[data-id="${CommentId}"][data-action="like"]`).classList.add('liked-btn');
                } else if (data.interaction === -1) {
                    document.querySelector(`.likedislikebtn[data-id="${CommentId}"][data-action="dislike"]`).classList.add('liked-btn');
                }
            });

        })
        .catch(err => {
            console.log(err);
        })
}