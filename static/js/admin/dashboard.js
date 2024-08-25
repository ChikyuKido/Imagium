let currentRoles = "";
let currentUserId = null;

let users = []

document.addEventListener('DOMContentLoaded',  () => {
     loadUsers();

    document.getElementById('users-tab').addEventListener('click', () => {
        showSection('users');
    });

    document.getElementById('settings-tab').addEventListener('click', () => {
        showSection('settings');
    });
});

function loadUsers(){
    const tbody = document.getElementById('users-table-body');
    tbody.innerHTML = '';
    fetch("/api/v1/admin/users")
         .then(response => response.json())
         .then(data => {
             users = data.users;
            users.forEach(user => {
                const tr = document.createElement('tr');
                tr.innerHTML = `
                    <td>${user.id}</td>
                    <td>${user.username}</td>
                    <td><button class="button is-small is-link" onclick="manageRoles(${user.id}, '${user.roles}')">${user.roles}</button></td>`;
                tbody.appendChild(tr);
            });
        })

}

function manageRoles(userId, roles) {
    currentRoles = roles;
    currentUserId = userId;
    document.getElementById('role-input').value = '';
    document.getElementById('role-modal').classList.add('is-active');
}

function closeModal() {
    document.getElementById('role-modal').classList.remove('is-active');
}

function deleteRole() {
    const role = document.getElementById('role-input').value;
    if (role && currentRoles.includes(role)) {
        currentRoles = currentRoles.split(',').filter(r => r !== role).join(',');
        updateRoles(currentUserId, currentRoles);
        closeModal();
    }
}

function addRole() {
    const role = document.getElementById('role-input').value;
    if (role && !currentRoles.includes(role)) {
        currentRoles += (currentRoles ? ',' : '') + role;
        updateRoles(currentUserId, currentRoles);
        closeModal();
    }
}

function updateRoles(userId, roles) {
    const user = users.find(user => user.id === userId);
    if (user) {
        user.roles = roles;
        fetch(`/api/v1/admin/users/changeRole/${userId}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ roles })
        }).then(value => loadUsers());
    }
}

function showSection(section) {
    document.getElementById('users-section').classList.add('is-hidden');
    document.getElementById('settings-section').classList.add('is-hidden');
    if (section === 'users') {
        document.getElementById('users-section').classList.remove('is-hidden');
        document.getElementById('users-tab').classList.add('is-active');
        document.getElementById('settings-tab').classList.remove('is-active');
    } else if (section === 'settings') {
        document.getElementById('settings-section').classList.remove('is-hidden');
        document.getElementById('settings-tab').classList.add('is-active');
        document.getElementById('users-tab').classList.remove('is-active');
    }
}