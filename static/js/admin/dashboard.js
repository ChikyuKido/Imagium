document.addEventListener('DOMContentLoaded', () => {
    loadUsers();
    loadSettings();

    document.getElementById('users-tab').addEventListener('click', () => {
        showSection('users');
    });

    document.getElementById('settings-tab').addEventListener('click', () => {
        showSection('settings');
    });
});

function loadUsers() {
    const tbody = document.getElementById('users-table-body');
    tbody.innerHTML = '';
    fetch("/api/v1/admin/users")
        .then(response => response.json())
        .then(data => {
            users = data.users;
            users.forEach(user => {
                const tr = document.createElement('tr');
                tr.innerHTML = `
                    <td>${user.ID}</td>
                    <td>${user.Username}</td>
                    <td><button class="button is-small is-link" onclick="manageRoles(${user.ID}, '${user.Roles}')">${user.Roles}</button></td>
                `;
                tbody.appendChild(tr);
            });
        });
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
    const user = users.find(user => user.ID === userId);
    if (user) {
        user.Roles = roles;
        fetch(`/api/v1/admin/users/changeRole/${userId}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ roles })
        }).then(() => loadUsers());
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

function loadSettings() {
    fetch('/api/v1/admin/settings')
        .then(response => response.json())
        .then(settings => {
            document.getElementById('admin-register').checked = settings.AdminRegister;
            document.getElementById('aggregation-time').value = settings.AggregationTime;
            document.getElementById('aggregation-job').value = settings.AggregationJob;
            document.getElementById('automatic-deletion-time').value = settings.AutomaticallyDeletionTime;
        });
}

function saveSettings() {
    const settings = {
        AdminRegister: document.getElementById('admin-register').checked,
        AggregationTime: parseInt(document.getElementById('aggregation-time').value, 10),
        AutomaticallyDeletionTime: parseInt(document.getElementById('automatic-deletion-time').value, 10)
    };

    fetch('/api/v1/admin/updateSettings', {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(settings)
    }).then(response => {
        if (response.ok) {
            alert('Settings updated successfully!');
        } else {
            alert('Failed to update settings.');
        }
    });
}
