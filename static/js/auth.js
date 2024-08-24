document.addEventListener('DOMContentLoaded', async () => {
    let neededPermission = document.getElementById("neededPermission").innerHTML;
    let redirectIfNotPermitted = document.getElementById("redirectIfNotPermitted").innerHTML;

    try {
        let token = localStorage.getItem('jwtToken');
        if(token == null) {
            token = "guest";
        }
        const response = await fetch('/api/v1/user/checkPermission', {
            method: 'POST',
            headers: {
                'Authorization': `${token}`,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ permission: neededPermission })
        });
        const data = await response.json();
        if (!data.hasPermission) {
            window.location = redirectIfNotPermitted
        }
    } catch (error) {
        console.error('Error checking permissions:', error);
    }
})