const app = new Vue({
    el: '#app',
    data: {
        config: {
            blocked_paths: '',
            allowed_ips: '',
            trusted_proxies: '',
            username: '',
            password: ''
        },
        shodanKeys: [],
        newShodanKey: '',
        logoutTimer: null,
        activeTab: 'about' // 默认激活 About 标签页
    },
    mounted() {
        this.loadConfig();
        this.loadShodanKeys();
        this.startLogoutTimer();
    },
    methods: {
        loadConfig() {
            axios.get('/get-config')
                .then(response => {
                    const data = response.data;
                    this.config = {
                        blocked_paths: Array.isArray(data.blocked_paths) ? data.blocked_paths.join('\n') : '',
                        allowed_ips: Array.isArray(data.allowed_ips) ? data.allowed_ips.join('\n') : '',
                        trusted_proxies: Array.isArray(data.trusted_proxies) ? data.trusted_proxies.join('\n') : '',
                        username: data.admin_user ? data.admin_user.username : '',
                        password: ''
                    };
                })
                .catch(error => {
                    console.error('Failed to load configuration:', error);
                    alert('Failed to load configuration: ' + error.response?.data || error.message);
                });
        },
        loadShodanKeys() {
            axios.get('/api/shodan-keys')
                .then(response => {
                    this.shodanKeys = response.data;
                })
                .catch(error => {
                    console.error('Failed to load Shodan keys:', error);
                    alert('Failed to load Shodan keys: ' + error.response?.data || error.message);
                });
        },
        addShodanKey() {
            if (this.newShodanKey.trim()) {
                const payload = JSON.stringify({ key: this.newShodanKey.trim() });
                axios.post('/api/shodan-keys', payload, {
                    headers: {
                        'Content-Type': 'application/json'
                    }
                })
                .then(() => {
                    this.loadShodanKeys();
                    this.newShodanKey = '';
                })
                .catch(error => {
                    console.error('Failed to add Shodan key:', error);
                    alert('Failed to add Shodan key: ' + (error.response?.data || error.message));
                });
            }
        },
        deleteShodanKey(key) {
            axios.delete('/api/shodan-keys', {
                data: { key: key },
                headers: {
                    'Content-Type': 'application/json'
                }
            })
            .then(response => {
                console.log('Delete response:', response);
                this.loadShodanKeys();
            })
            .catch(error => {
                console.error('Failed to delete Shodan key:', error);
                const errorMessage = error.response?.data || error.message;
                alert('Failed to delete Shodan key: ' + errorMessage);
            });
        },
        updateConfig() {
            const configToSave = {
                blocked_paths: this.getTextareaValues(this.config.blocked_paths),
                allowed_ips: this.getTextareaValues(this.config.allowed_ips),
                trusted_proxies: this.getTextareaValues(this.config.trusted_proxies),
                admin_user: {
                    username: this.config.username,
                    ...(this.config.password && { password: this.config.password })
                }
            };

            axios.post('/update-config', configToSave)
                .then(() => {
                    alert('Configuration updated successfully');
                    this.loadConfig();
                })
                .catch(error => {
                    console.error('Error updating configuration:', error);
                    alert('An error occurred while updating the configuration: ' + (error.response?.data || error.message));
                });
            this.resetLogoutTimer();
        },
        logout() {
            axios.post('/logout')
                .then(() => {
                    window.location.href = '/login';
                })
                .catch(error => {
                    alert('Failed to logout: ' + (error.response?.data || error.message));
                });
        },
        startLogoutTimer() {
            if (this.logoutTimer) {
                clearTimeout(this.logoutTimer);
            }
            this.logoutTimer = setTimeout(() => {
                this.logout();
            }, 10 * 60 * 1000);
        },
        resetLogoutTimer() {
            this.startLogoutTimer();
        },
        getTextareaValues(text) {
            return text.split('\n').filter(item => item.trim() !== '');
        }
    }
});
