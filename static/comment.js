let offset = 3

commentForm.addEventListener('submit', function (e) {
    e.preventDefault();

    const formData = new FormData(this);

    fetch('/Commenting', {
        method: 'POST',
        body: formData
    }).then(resp => {
        if (resp.status === 401) {
            window.location.href = '/login';
            return;
        }

        if (resp.status === 400) {
            return resp.json();
        }

        if (!resp.ok) {
            window.location.replace(`/error?s=${resp.status}&st=${resp.statusText}`)
            return
        }
        return resp.json();
    })
        .then(comment => {
            if (comment.error) {
                const errorMsg = document.createElement('div');
                errorMsg.className = "error-message";
                errorMsg.style.color = 'red';
                errorMsg.textContent = comment.error;

                const existingError = commentForm.querySelector('.error-message');
                if (existingError) {
                    console.log(existingError);
                    existingError.remove();
                }

                commentForm.insertBefore(errorMsg, commentForm.firstChild);
                return
            }
            const commentElement = createCommentElement(comment);
            commentsContainer.insertBefore(commentElement, commentsContainer.firstChild);
            offset += 1
            document.querySelector('.myform').reset();
        }).catch(error => {
            console.error("error", error);
        });
})

function createCommentElement(comment) {
    const div = document.createElement('div');
    div.className = "comment"
    div.innerHTML = `
        <h3>${comment.Uname}</h3>
        <p>${comment.Content}</p>

        <button
        class="action-btn likedislikebtn"
        data-id="${comment.Id}"
        data-action="like"
        data-type="comment"
        >
        <i class="fas fa-thumbs-up"></i>
        <span id="like_post-${comment.Id}">0</span>
        </button>
        <button
        class="action-btn likedislikebtn"
        data-id="${comment.Id}"
        data-action="dislike"
        data-type="comment"
        >
        <i class="fas fa-thumbs-down"></i>
        <span id="dislike_post-${comment.Id}">0</span>
        </button>
    `;
    return div;
}


// load more comments

let loading = false
let noMoreComments = false

function loadMorecomments() {
    if (loading || noMoreComments) return;

    loading = true;
    loadingContainer.style.visibility = "visible"

    const post_id = document.querySelector('input[name="post_id"]').value;

    const url = new URL('/load-more-comments', window.location.origin);
    url.searchParams.set('offset', offset);
    url.searchParams.set('post_id', post_id); 

    fetch(url)
        .then(resp => resp.json())
        .then(comments => {
            if (comments == null) {
                noMoreComments = true
                loadingContainer.textContent = 'No more comments to load';
                return;
            }
            console.log(comments);
            
            comments.forEach(post => {
                const postElement = createCommentElement(post);
                commentsContainer.appendChild(postElement);
            })

            offset += comments.length
            loading = false;
            loadingContainer.style.visibility = "hidden"
        })
        .catch(err => {
            console.error('Error loading more comments:', err);
            loading = false;
            loadingContainer.style.visibility = "hidden"
        })
}

const observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
        if (entry.isIntersecting) {
            loadMorecomments()
        }
    })
})

observer.observe(loadingContainer)