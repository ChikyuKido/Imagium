const apiUrl = '/api/v1/user/library';
let currentPage = 0;

const fetchImages = async (page) => {
    try {
        const response = await fetch(`${apiUrl}?site=${currentPage}`);
        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Error fetching images:', error);
    }
};

const renderImages = (images) => {
    const imageCardsContainer = document.getElementById('image-cards');
    imageCardsContainer.innerHTML = '';

    images.forEach(image => {
        const card = `
                    <a class="column is-one-quarter" href="/image/view/${image.UUID}">
                        <div class="card">
                            <div class="card-image">
                                <figure class="image">
                                    <img src="/image/get/${image.UUID}.webp?resize=800x600&quality=80" alt="${image.Name}">
                                </figure>
                            </div>
                            <div class="card-content">
                                <div class="media">
                                    <div class="media-content">
                                        <p class="title is-4">${image.Name}</p>
                                        <p class="is-6">Size: ${image.Size} bytes</p>
                                        <p class="is-6">SubImages: ${image.SubImages}</p>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </a>
                `;
        imageCardsContainer.innerHTML += card;
    });
};

const renderPagination = (totalPages) => {
    const paginationContainer = document.getElementById('pagination');
    paginationContainer.innerHTML = '';

    for (let i = 0; i <= totalPages; i++) {
        const pageLink = `
                    <a class="pagination-link ${i === currentPage ? 'is-current' : ''}" href="?site=${i}" data-page="${i}">
                        ${i+1}
                    </a>
                `;
        paginationContainer.innerHTML += pageLink;
    }
};

const loadPage = async (page) => {
    const images = await fetchImages(page);
    if (images) {
        renderImages(images["Images"]);
        renderPagination(images["Pages"]);
    }
};

document.addEventListener('DOMContentLoaded', () => {
    const urlParams = new URLSearchParams(window.location.search);
    currentPage = parseInt(urlParams.get('site')) || 0;
    loadPage(currentPage);

    document.getElementById('pagination').addEventListener('click', (event) => {
        if (event.target.matches('.pagination-link')) {
            event.preventDefault();
            const newPage = parseInt(event.target.dataset.page);
            if (newPage !== currentPage) {
                currentPage = newPage;
                loadPage(currentPage);
            }
        }
    });
});