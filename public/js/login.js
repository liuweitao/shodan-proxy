new Vue({
    el: '#app',
    data: {
        username: '',
        password: ''
    },
    methods: {
        submitForm() {
            if (!this.username || !this.password) {
                alert('请输入用户名和密码。');
                return;
            }
            
            axios.post('/login', {
                username: this.username,
                password: this.password
            })
            .then(response => {
                console.log('Login successful, redirecting to admin page');
                window.location.href = '/admin';
            })
            .catch(error => {
                console.error('Login error:', error);
                alert('登录失败: ' + (error.response?.data || error.message));
            });
        }
    }
});
