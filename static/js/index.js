document.addEventListener('DOMContentLoaded', () => {
    const uploadArea = document.getElementById('uploadArea');
    const fileInput = document.getElementById('fileInput');
    const spinner = document.getElementById('spinner');
    const helpText = document.getElementById('help-text');

    uploadArea.addEventListener('dragover', (event) => {
        event.preventDefault();
        uploadArea.classList.add('dragover');
    });

    uploadArea.addEventListener('dragleave', () => {
        uploadArea.classList.remove('dragover');
    });

    uploadArea.addEventListener('drop', (event) => {
        event.preventDefault();
        uploadArea.classList.remove('dragover');
        handleFile(event.dataTransfer.files);
    });

    uploadArea.addEventListener('click', () => {
        fileInput.click();
    });

    fileInput.addEventListener('change', () => {
        handleFile(fileInput.files);
    });
    function handleFile(files)  {
        if (files.length > 0) {
            const file = files[0];
            const formData = new FormData();
            formData.append('file', file);

            spinner.style.visibility = 'visible';
            helpText.style.visibility = 'hidden';

            fetch('/api/v1/image/uploadImage', {
                method: 'POST',
                body: formData
            })
                .then(async response => {
                    let data = await response.json();
                    if (response.status !== 200) {
                        alert('Upload error: ' + data.error)
                    } else {
                        window.location = data.url
                    }
                })
                .catch(error => alert('Upload error:' + error))
                .finally(() => {
                    spinner.style.visibility = 'hidden';
                    helpText.style.visibility = 'visible';
                });
        } else {
            console.log('No file selected.');
        }
    }
});
