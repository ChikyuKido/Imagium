document.addEventListener('DOMContentLoaded', () => {
    const form = document.getElementById('loginForm') || document.getElementById('registerForm') || document.getElementById('adminRegisterForm');

    if (!form) {
        console.error('Form not found');
        return;
    }

    form.addEventListener('submit', async (event) => {
        event.preventDefault();

        const formData = new FormData(form);
        const data = {
            username: formData.get('username'),
            password: formData.get('password')
        };

        let url = '';
        let redirectUrl = '';

        switch (form.id) {
            case 'loginForm':
                url = '/api/v1/user/login';
                redirectUrl = '/';
                break;
            case 'registerForm':
                url = '/api/v1/user/register';
                redirectUrl = '/login';
                break;
            case 'adminRegisterForm':
                url = '/api/v1/admin/register';
                redirectUrl = '/login';
                break;
            default:
                console.error('Unknown form id');
                return;
        }

        try {
            const response = await fetch(url, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(data)
            });

            if (response.ok) {
                window.location.href = redirectUrl;
            } else {
                const error = await response.json();
                console.error('Operation failed', error);
                alert('Operation failed: ' + error.message);
            }
        } catch (error) {
            console.error('Error during operation', error);
            alert('An error occurred. Please try again.');
        }
    });
});
