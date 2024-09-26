document.addEventListener('DOMContentLoaded', function() {
    console.log('DOM fully loaded and parsed');
    initializeForms();
    initializePageBasedFunctions();
    setupLogout();
});

function initializeForms() {
    setupFormSubmission('adminSignupForm', handleAdminSignup);
    setupFormSubmission('adminLoginForm', handleAdminLogin);
    setupFormSubmission('userSignupForm', handleUserSignup);
    setupFormSubmission('userLoginForm', handleUserLogin);
    setupFormSubmission('createUserForm', handleCreateUser);
    setupFormSubmission('updateUserForm', handleUpdateUser);
    setupFormSubmission('deleteUserForm', handleDeleteUser);

    document.getElementById('getUsersButton')?.addEventListener('click', fetchUsers);
}

function setupFormSubmission(formId, handler) {
    const form = document.getElementById(formId);
    if (form) {
        form.addEventListener('submit', handler);
    }
}

document.getElementById('adminSignupForm').addEventListener('submit', handleAdminSignup);

async function handleAdminSignup(event) {
    event.preventDefault();
    const formData = new FormData(event.target);
    const data = extractFormData(formData, ['first_name', 'last_name', 'email', 'password', 'phone']);
    await sendRequest('http://localhost:8000/admin/signup', 'POST', data, handleAdminSignupResponse);
}

function handleAdminSignupResponse(response, result) {
    console.log("Response status:", response.status);
    console.log("Result:", result);
    if (response.ok) {
        alert('Admin signup successful');
        // Redirect to admin login after successful signup
        window.location.href = '/frontend/admin_login.html'; 
    } else {
        const errorMessage = result && result.error ? result.error : 'Signup failed, try using another email';
        // Debug log
        console.log("Error message:", errorMessage); 
        alert('Error: ' + errorMessage);
    }
}

function extractFormData(formData, fields) {
    const data = {};
    fields.forEach(field => {
        data[field] = formData.get(field);
    });
    return data;
}

async function sendRequest(url, method, data, callback) {
    try {
        const response = await fetch(url, {
            method: method,
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        });
        // Log status code
        console.log("Response status:", response.status); 

        let result;
        try {
            result = await response.json();
        } catch (error) {
            console.error("Failed to parse JSON:", error);
            result = {error: "Try after 1 hour" };
        }
        console.log("Parsed result:", result);
        callback(response, result);
    } catch (error) {
        console.error("Network error:", error);
        alert('Network error: ' + error.message);
    }
}



async function handleAdminLogin(event) {
    event.preventDefault();
    const formData = new FormData(event.target);
    const data = extractFormData(formData, ['email', 'password']);
    await sendRequest('http://localhost:8000/admin/login', 'POST', data, handleAdminLoginResponse);
}

async function handleUserSignup(event) {
    event.preventDefault();
    const formData = new FormData(event.target);
    const data = extractFormData(formData, ['first_name', 'last_name', 'email', 'password', 'phone', 'user_type']);
    await sendRequest('http://localhost:8000/users/signup', 'POST', data, handleUserSignupResponse);
}

function handleUserSignupResponse(response, result) {
    console.log("Response status:", response.status);
    console.log("Response received:", response);
    console.log("Result:", result);
   
    if (response.ok) {
        alert('User signup successful');
        window.location.href = '/frontend/user_login.html'; // Redirect to user login after successful signup
    } else {
        const errorMessage = result && result.error ? result.error : 'Signup failed, try using another email';
        // Debug log
        console.log("Error message:", errorMessage); 
        alert('Error: ' + errorMessage);
    }
}

function extractFormData(formData, fields) {
    const data = {};
    fields.forEach(field => {
        data[field] = formData.get(field);
    });
    return data;
}

async function sendRequest(url, method, data, callback) {
    try {
        const response = await fetch(url, {
            method: method,
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        });
        // Log status code
        console.log("Response status:", response.status); 

        let result;
        try {
            result = await response.json();
        } catch (error) {
            console.error("Failed to parse JSON:", error);
            result = {error: "Try after 1 hour" };
        }
        console.log("Parsed result:", result);
        callback(response, result);
    } catch (error) {
        console.error("Network error:", error);
        alert('Network error: ' + error.message);
    }
}

document.getElementById('userSignupForm').addEventListener('submit', handleUserSignup);


async function handleUserLogin(event) {
    event.preventDefault();
    const formData = new FormData(event.target);
    const data = extractFormData(formData, ['email', 'password']);
   // const email = data.email; 

    console.log('Sending login request with data:', data); // Debug log
    await sendRequest('http://localhost:8000/users/login', 'POST', data, handleUserLoginResponse);
}

function extractFormData(formData, keys) {
    const data = {};
    keys.forEach(key => {
        data[key] = formData.get(key);
    });
    return data;
}

async function sendRequest(url, method, data, callback) {
    console.log('Data being sent:', data); // Debug log
    const response = await fetch(url, {
        method: method,
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    });
    const result = await response.json();
    console.log('Received response:', result); // Debug log
    callback(response, result, data.email);
}

function handleUserLoginResponse(response, result, email) {
    console.log('Email passed to handleUserLoginResponse:', email); // Debug log

    if (response.ok) {
        console.log('Login successful, storing token and username'); // Debug log
        localStorage.setItem('userToken', result.token);
        localStorage.setItem('username', result.username);  // Store the username
        console.log('Stored username:', email); 
        window.location.href = '/frontend/user_home.html';
    } else {
        console.error('Login failed:', result.error); // Debug log
        alert('Error: ' + result.error);  // Display error message from the backend
    }
}

document.getElementById('userLoginForm').addEventListener('submit', handleUserLogin);



async function handleCreateUser(event) {
    event.preventDefault();
    const formData = new FormData(event.target);
    const data = extractFormData(formData, ['first_name', 'last_name', 'email', 'password', 'phone', 'user_type']);
    await sendRequest('http://localhost:8000/admin/create_user', 'POST', data, handleCreateUserResponse);
}

function handleCreateUserResponse(response, result) {
    console.log("Response status:", response.status);
    console.log("Response received:", response);
    console.log("Result:", result);
    if (response.ok) {
        alert('User created successfully');
        // Refresh the user list here
        fetchUsers(); 
    } else {
        const errorMessage = result && result.error ? result.error : 'User creation failed, try using another email';
        console.log("Error message:", errorMessage);
        alert('Error: ' + errorMessage);
    }
}

document.getElementById('createUserForm').addEventListener('submit', handleCreateUser);


document.getElementById('updateUserForm').addEventListener('submit', handleUpdateUser);

async function handleUpdateUser(event) {
    event.preventDefault();
    const formData = new FormData(event.target);
    const data = extractFormData(formData, ['userId', 'firstName', 'lastName', 'email', 'phone', 'userType']);
    const userId = data.userId;
    delete data.userId; // Remove userId from the data to be sent in the body
    await sendRequest(`http://localhost:8000/admin/update_user/${userId}`, 'PUT', data, handleUpdateUserResponse);
}

function handleUpdateUserResponse(response, result) {
    if (response.ok) {
        alert('User updated successfully');
        // Optionally refresh user data or redirect
        fetchUsers();
    } else {
        const errorMessage = result && result.error ? result.error : 'Update failed, email already exist';
        // Debug log
        console.log("Error message:", errorMessage); 
        alert('Error: ' + errorMessage);
        response.json().then(errorData => {
            alert(`Error: ${errorData.error}`);
        }).catch(() => {
            alert('Error'+ error.message);
        });
    }
}

function extractFormData(formData, fields) {
    const data = {};
    fields.forEach(field => {
        data[field] = formData.get(field);
    });
    return data;
}

async function sendRequest(url, method, data, callback) {
    try {
        const response = await fetch(url, {
            method: method,
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + localStorage.getItem('adminToken')
            },
            body: JSON.stringify(data)
        });
        const result = await response.json();
        callback(response, result);
    } catch (error) {
        console.error('Error:', error);
    }
}

function handleAdminLoginResponse(response, result) {
    if (response.ok) {
        console.log('Storing admin token:', result.token);
        localStorage.setItem('adminToken', result.token);
        alert('Logged in successfully');
        window.location.href = '/frontend/dashboard.html';
    } else {
        alert('Error: Invalid username or password');
    }
}

function handleUserLoginResponse(response, result, email) {
    if (response.ok) {
        localStorage.setItem('userToken', result.token);
        localStorage.setItem('username', result.username);
        window.location.href = '/frontend/user_home.html';
    } else {
        alert('Error: Invalid username or password, please contact admin');
    }
}

document.getElementById('userLoginForm').addEventListener('submit', handleUserLogin);

// Initialize page-based functions
function initializePageBasedFunctions() {
    ensureAuthentication();

    if (window.location.pathname.endsWith('/frontend/user_home.html')) {
        fetchUserDetails();
    } else if (window.location.pathname.endsWith('/frontend/dashboard.html')) {
        fetchUsers();
    }
}

// Ensure authentication on page load
document.addEventListener('DOMContentLoaded', (event) => {
    ensureAuthentication();
    initializePageBasedFunctions();
});

// Handle back and forward navigation
window.addEventListener('popstate', (event) => {
    ensureAuthentication();
    initializePageBasedFunctions();
});






















        