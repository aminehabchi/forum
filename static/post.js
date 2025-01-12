let offset = 0;
let loading = false;
let noMorePosts = false;
let currentFilter = "";

function loadMorePosts() {
  if (loading || noMorePosts) return;

  loading = true;
  loadingContainer.style.visibility = "visible";

  const url = new URL("/load-more-posts", window.location.origin);
  url.searchParams.set("offset", offset);

  if (currentFilter) {
    url.searchParams.set("type", currentFilter);
  }

  fetch(url)
    .then(async (resp) => {
      if (resp.ok) {
        return resp.json();
      } else {
        const html = await resp.text();
        document.documentElement.innerHTML = html;
      }
    })
    .then((posts) => {
      if (posts == null) {
        noMorePosts = true;
        loadingContainer.textContent = "No more posts to load";
        return;
      }

      posts.forEach((post) => {
        const postElement = createPostElement(post);
        postsContainer.appendChild(postElement);
      });

      offset += posts.length;
      loading = false;
      loadingContainer.style.visibility = "hidden";
    })
    .catch((err) => {
      console.error("Error loading more posts:", err);
      loading = false;
      loadingContainer.style.visibility = "hidden";
    });
}

function createPostElement(post) {
  let img = "";
  if (post.ImgBase64!= undefined && post.ImgBase64!="") {
    img = `<img id="image" alt="Image" src="data:image/png;base64,${post.ImgBase64}" />`;
  }
  const div = document.createElement("div");
  div.className = "post";
  div.innerHTML = `
        <div class="profilInfo">
            <img class="profileImg" src="static/profil.png">
            <div class="profile-details">
                <span>${post.Name}</span>
                <span class="time">${post.CreatedAt}</span>
            </div>
        </div>
        <h4>${post.Title}</h4>
            ${post.Category.map(
              (cat) => `<a class="category" href="#"
                onclick="handleCategoryClick(event, '${cat}')"
              ><i class="fas fa-tag"></i> ${cat}</a
            >`
            ).join(" ")}
        <p class="content">${post.Content}</p>
       ${img}
        <div style="margin-top: 15px;">
            <button class="action-btn likedislikebtn ${
              post.UserInteraction === 1 ? "liked-btn" : ""
            }"
                    data-id="${post.ID}" data-action="like" data-type="post">
                <i class="fas fa-thumbs-up"></i>
                <span id="like_post-${post.ID}">${post.Likes}</span>
            </button>
            <button class="action-btn likedislikebtn ${
              post.UserInteraction === -1 ? "liked-btn" : ""
            }"
                    data-id="${post.ID}" data-action="dislike" data-type="post">
                <i class="fas fa-thumbs-down"></i>
                <span id="dislike_post-${post.ID}">${post.Dislikes}</span>
            </button>
            <a class="action-btn" href="/Commenting?post_id=${post.ID}">
                <i class="fas fa-comment"></i>${post.NbComment}
            </a>
        </div>
    `;
  return div;
}

const observer = new IntersectionObserver((entries) => {
  entries.forEach((entry) => {
    if (entry.isIntersecting) {
      loadMorePosts();
    }
  });
});

observer.observe(loadingContainer);

// ------ Filter ------- //

document.addEventListener("DOMContentLoaded", () => {
  document.querySelectorAll(".filteraction").forEach((r) => {
    r.checked = false;
    r.setAttribute("data-checked", "false");
  });
});

filterContainer.addEventListener("click", (e) => {
  const action = e.target.closest(".filteraction");
  if (!action) return;

  if (action.getAttribute("data-checked") === "true") {
    action.checked = false;
    action.setAttribute("data-checked", "false");
    filterCategory("");
  } else {
    document.querySelectorAll(".filteraction").forEach((r) => {
      r.checked = false;
      r.setAttribute("data-checked", "false");
    });

    action.checked = true;
    action.setAttribute("data-checked", "true");

    filterCategory(action.value);
  }
});

function filterCategory(type) {
  offset = 4;
  noMorePosts = false;
  currentFilter = type;

  fetch(`/filter?type=${type}`)
    .then((resp) => {
      if (resp.status === 401) {
        window.location.href = "/login";
        return;
      }

      if (!resp.ok) {
        window.location.replace(
          `/error?s=${resp.status}&st=${resp.statusText}`
        );
        return;
      }
      return resp.json();
    })
    .then((posts) => {
      postsContainer.innerHTML = "";
      if (posts == null) {
        return;
      }
      posts.forEach((post) => {
        const postElement = createPostElement(post);
        postsContainer.appendChild(postElement);
      });
    })
    .catch((error) => console.error("Error:", error));
}

function handleCategoryClick(event, category) {
  event.preventDefault();

  const action = document.querySelector(`.filteraction[value=${category}]`);
  if (action) {
    document.querySelectorAll(".filteraction").forEach((r) => {
      r.checked = false;
      r.setAttribute("data-checked", "false");
    });

    action.checked = true;
    action.setAttribute("data-checked", "true");
  }
  filterCategory(category);
}
