document.addEventListener('DOMContentLoaded', () => {
    const form = document.getElementById('registerForm');
    form.addEventListener('submit', async (event) => {
        event.preventDefault();

        const formData = new FormData(form);
        const data = {
            username: formData.get('username'),
            password: formData.get('password')
        };

        try {
            const response = await fetch('/api/v1/admin/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(data)
            });

            if (response.ok) {
                const result = await response.json();
                window.location.href = '/login';
            } else {
                const error = await response.json();
                console.error('Register failed', error.error);
                alert('Register failed: ' + error.error);
            }
        } catch (error) {
            console.error('Error during Register', error);
            alert('An error occurred. Please try again.');
        }
    });
});