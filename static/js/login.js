document.addEventListener('DOMContentLoaded', () => {
    const form = document.getElementById('loginForm');
    form.addEventListener('submit', async (event) => {
        event.preventDefault();

        const formData = new FormData(form);
        const data = {
            username: formData.get('username'),
            password: formData.get('password')
        };

        try {
            const response = await fetch('/api/v1/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(data)
            });

            if (response.ok) {
                const result = await response.json();
                const token = result.token;
                if (token) {
                    localStorage.setItem('jwt', token);
                    window.location.href = '/';
                } else {
                    console.error('JWT not received');
                    alert('Login failed: No JWT received');
                }
            } else {
                const error = await response.json();
                console.error('Login failed', error);
                alert('Login failed: ' + error.message);
            }
        } catch (error) {
            console.error('Error during login', error);
            alert('An error occurred. Please try again.');
        }
    });
});