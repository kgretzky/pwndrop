var appHome = Vue.component("app-home", {
    template: `<div id="app">
    <div v-if="!isDead && isLoggedIn" class="top-left-bar">
        <button class="btn btn-primary btn-circle" @click="showConfig()">
            <i class="fas fa-cog"></i>
        </button>
    </div>
    <div v-if="!isDead" class="top-right-bar">
        <button class="btn btn-primary btn-circle" @click="logout()">
            <i class="fas fa-sign-out-alt"></i>
        </button>
    </div>

    <div class="bg-title">
        <a href="/"><img src="pages/img/pwndrop-title.png" alt="pwndrop title" /></a>
    </div>
    <div class="bg-footer">
        made by <a href="https://twitter.com/mrgretzky" target="_blank">@mrgretzky</a>
    </div>
    <div class="bg-version">
        version {{ version }}
    </div>
    <b-modal v-model="configShow" id="config-modal" title="Settings" size="lg" ok-title="Save"
        @ok.prevent="saveConfig()"
        :disabled="errors.any() || !isConfigComplete"
        >
        <form>
            <div class="form-group row">
                <label for="redirect-url" class="col-sm-3 col-form-label label-help">Redirect URL:
                    <i class="fas fa-question-circle label-qmark" v-tooltip:bottom="'Visitors will be redirected to this URL if they provide a wrong download URL or are unauthorized to view the admin panel'"></i>
                </label>
                <div class="col-sm-9">
                    <input type="text" class="form-control" id="redirect-url" spellcheck="false"
                        v-model="config.redirect_url"
                        name="redirect-url"
                        v-bind:class="{'form-control': true, 'error': errors.has('redirect-url') }"
                    >
                    <div v-show="errors.has('redirect-url')" class="form-error">{{ errors.first('redirect-url') }}</div>
                </div>
            </div>
            <div class="form-group row">
                <label for="secret-path" class="col-sm-3 col-form-label label-help">Secret Path:
                    <i class="fas fa-question-circle label-qmark" v-tooltip:bottom="'Visiting this path in a browser will authorize the visitor to view the admin panel (IMPORTANT! CHANGE FROM DEFAULT)'"></i>
                </label>
                <div class="col-sm-9">
                    <input type="text" class="form-control" id="secret-path" spellcheck="false"
                        v-model="config.secret_path"
                        name="secret-path"
                        v-validate="'required'"
                        v-bind:class="{'form-control': true, 'error': errors.has('secret-path') }"
                    >
                    <div v-show="errors.has('secret-path')" class="form-error">{{ errors.first('secret-path') }}</div>
                </div>
            </div>
            <div class="form-group row">
                <label for="cookie-name" class="col-sm-3 col-form-label label-help">Secret-Cookie Name:
                    <i class="fas fa-question-circle label-qmark" v-tooltip:bottom="'Secret cookie name, which is used for authorizing the visitor to view the admin panel'"></i>
                </label>
                <div class="col-sm-9">
                    <input type="text" class="form-control" id="cookie-name" spellcheck="false"
                        v-model="config.cookie_name"
                        name="cookie-name"
                        v-validate="'required'"
                        v-bind:class="{'form-control': true, 'error': errors.has('cookie-name') }"
                    >
                    <div v-show="errors.has('cookie-name')" class="form-error">{{ errors.first('cookie-name') }}</div>
                </div>
            </div>
            <div class="form-group row">
                <label for="cookie-token" class="col-sm-3 col-form-label label-help">Secret-Cookie Value:
                    <i class="fas fa-question-circle label-qmark" v-tooltip:bottom="'Secret cookie value, which is used for authorizing the visitor to view the admin panel'"></i>
                </label>
                <div class="col-sm-9">
                    <input type="text" class="form-control" id="cookie-token" spellcheck="false"
                        v-model="config.cookie_token"
                        name="cookie-token"
                        v-validate="'required'"
                        v-bind:class="{'form-control': true, 'error': errors.has('cookie-token') }"
                    >
                    <div v-show="errors.has('cookie-token')" class="form-error">{{ errors.first('cookie-token') }}</div>
                </div>
            </div>
        </form>
    </b-modal>

    <div class="bg-logo"></div>
    <div class="container">
        <div v-if="isDead" class="text-center">
            <span class="big-icon">
                <i class="fas fa-dizzy"></i>
            </span>
        </div>
        <div v-else-if="!isLoaded"></div>
        <div v-else>
            <router-view></router-view>
        </div>
    </div>
</div>`,
    $_veeValidate: {
        validator: "new"
    },
	name: "app",
	data() {
		return {
			url: Config.Hostname + Config.AdminDir + "/" + Config.ApiPath,
			isLoaded: false,
			isLoggedIn: false,
			doCreateAccount: false,
			doLogin: false,
			isDead: false,
			Username: "",
			config: {
				secret_path: "",
				redirect_url: "",
				cookie_name: "",
				cookie_token: ""
            },
            configShow: false,
            version: "-"
		};
    },
    computed: {
        isConfigComplete () {
            return this.config.secret_path && this.config.cookie_name && this.config.cookie_token;
        }
    },
	methods: {
		authCheck() {
			axios
				.get(this.url + "/auth")
				.then(response => {
					console.log(response);
					if (response.data.data.status == 0) {
                        this.doCreateAccount = true;
                        this.$router.push('/create_account').catch(err => {});
						return;
					} else if (response.data.data.status == 1) {
						this.doCreateAccount = false;
						this.doLogin = false;
						this.isLoggedIn = true;
					}
				})
				.catch(error => {
					console.log(error);
					this.doLogin = true;
                    this.$router.push('/login').catch(err => {});
				})
				.then(() => {
					this.isLoaded = true;
				});
		},
		logout() {
			if (this.isLoggedIn) {
				axios
					.get(this.url + "/logout")
					.then(response => {
						console.log(response);
						this.doLogin = true;
                        this.isLoggedIn = false;
                        this.$router.push('/login').catch(err => {});
					})
					.catch(error => {
						console.log(error);
					});
			} else {
				axios
					.get(this.url + "/clear_secret")
					.then(response => {
						console.log(response);
						this.isDead = true;
					})
					.catch(error => {
						console.log(error);
					});
			}
		},
		showConfig() {
			axios
				.get(this.url + "/config")
				.then(response => {
					console.log(response);
					var r = response.data.data;
					this.config.secret_path = r.secret_path;
					this.config.redirect_url = r.redirect_url;
					this.config.cookie_name = r.cookie_name;
					this.config.cookie_token = r.cookie_token;

					this.$bvModal.show("config-modal");
				})
				.catch(error => {
					console.log(error);
				});
		},
		saveConfig() {
            if (!this.isConfigComplete) {
                return;
            }
			axios
				.post(
					this.url + "/config",
					{
						secret_path: this.config.secret_path,
						redirect_url: this.config.redirect_url,
						cookie_name: this.config.cookie_name,
						cookie_token: this.config.cookie_token
					},
					{
						headers: {
							"content-type": "application/json"
						}
					}
				)
				.then(response => {
					console.log(response);
					this.$bvModal.hide("config-modal");
				})
				.catch(error => {
					console.log(error);
				});
        },
        syncVersion() {
            axios
                .get(this.url + "/version")
                .then(response => {
                    var r = response.data.data;

                    this.version = r.version;
                })
                .catch(error => {
                    console.log(error);
                });
        }
    },
	created() {
        this.syncVersion();
        this.authCheck();
        
        this.mainBus.$on('createdAccount', () => {
			this.doCreateAccount = false;
			this.doLogin = true;
            this.$router.push('/login').catch(err => {});
        });
        this.mainBus.$on('loggedIn', (username) => {
			console.log(username);
			this.Username = username;
			this.doLogin = false;
			this.isLoggedIn = true;
            this.$router.push('/').catch(err => {});
        });
	}
});