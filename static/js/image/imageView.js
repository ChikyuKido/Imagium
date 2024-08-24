function updateURL() {
    const rootURL = `${window.location.protocol}//${window.location.host}`;
    const fullPath = window.location.pathname;
    const pathSegments = fullPath.split('/');
    const lastSegment = pathSegments[pathSegments.length - 1] || '/';

    let url = rootURL + '/image/get/' + lastSegment;

    const fileType = document.getElementById('file-type').value;
    url += `.${fileType}`;

    const resize = document.getElementById('resize').checked;
    if (resize) {
        const width = document.getElementById('resize-width').value;
        const height = document.getElementById('resize-height').value;
        url += `?resize=${width}x${height}`;
    }

    const quality = document.getElementById('quality').checked;
    if (quality) {
        const qualityValue = document.getElementById('quality-value').value;
        url += (url.includes('?') ? '&' : '?') + `quality=${qualityValue}`;
    }

    const crop = document.getElementById('crop').checked;
    if (crop) {
        const cropWidth = document.getElementById('crop-width').value;
        const cropHeight = document.getElementById('crop-height').value;
        const cropX = document.getElementById('crop-x').value;
        const cropY = document.getElementById('crop-y').value;
        url += (url.includes('?') ? '&' : '?') + `crop=${cropWidth}x${cropHeight}:${cropX}:${cropY}`;
    }

    const blur = document.getElementById('blur').checked;
    if (blur) {
        url += (url.includes('?') ? '&' : '?') + 'blur=true';
    }

    document.getElementById('url-preview').textContent = url;
    document.getElementById('html-code').textContent = `<img src="${url}" alt="Image">`;
}
function decodeHtmlEntities(html) {
    const textarea = document.createElement('textarea');
    textarea.innerHTML = html;
    return textarea.value;
}
function copyToClipboard(text) {
    navigator.clipboard.writeText(text).then(function() {
        console.log('Text copied to clipboard');
    }).catch(function(error) {
        console.error('Failed to copy text: ', error);
    });
}
document.addEventListener('DOMContentLoaded', () => {
    const rootURL = `${window.location.protocol}//${window.location.host}`;
    const fullPath = window.location.pathname;
    const pathSegments = fullPath.split('/');
    const lastSegment = pathSegments[pathSegments.length - 1] || '/';

    let url = rootURL + '/image/get/' + lastSegment+".png";
    document.getElementById("thumbnail").src = url;
    document.querySelectorAll('input, select').forEach(element => {
        element.addEventListener('change', updateURL);
    });
    updateURL();

    document.getElementById('resize').addEventListener('change', function() {
        document.getElementById('resize-options').classList.toggle('hidden', !this.checked);
        updateURL();
    });

    document.getElementById('quality').addEventListener('change', function() {
        document.getElementById('quality-options').classList.toggle('hidden', !this.checked);
        updateURL();
    });

    document.getElementById('crop').addEventListener('change', function() {
        document.getElementById('crop-options').classList.toggle('hidden', !this.checked);
        updateURL();
    });

    document.getElementById('blur').addEventListener('change', updateURL);

    document.getElementById('url-preview').addEventListener('click', function() {
        copyToClipboard( document.getElementById('url-preview').innerHTML)
    });
    document.getElementById('html-code').addEventListener('click', function() {
        const htmlContent = document.getElementById('html-code').innerHTML;
        const decodedContent = decodeHtmlEntities(htmlContent);
        copyToClipboard(decodedContent);
    });
});