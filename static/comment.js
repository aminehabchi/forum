commentForm.addEventListener('submit', function (e) {
    e.preventDefault();

    const formData = new FormData(this);

    fetch('/Commenting', {
        method: 'POST',
        body: formData
    }).then(resp => resp.json())
        .then(comment => {
            const commentElement = createCommentElement(comment);
            commentsContainer.insertBefore(commentElement, commentsContainer.firstChild);
            document.querySelector('.myform').reset();
        }).catch(error => {
            console.error(error);
        });
})

function createCommentElement(comment) {
    const div = document.createElement('div');
    div.className = "comment"
    div.innerHTML = `
        <h3>${comment.uname}</h3>
        <p>${comment.content}</p>

        <button
        class="action-btn likedislikebtn"
        data-id="${comment.id}"
        data-action="like"
        data-type="comment"
        >
        <i class="fas fa-thumbs-up"></i>
        <span id="like_post-${comment.id}">0</span>
        </button>
        <button
        class="action-btn likedislikebtn"
        data-id="${comment.id}"
        data-action="dislike"
        data-type="comment"
        >
        <i class="fas fa-thumbs-down"></i>
        <span id="dislike_post-${comment.id}">0</span>
        </button>
    `;
    return div;
}