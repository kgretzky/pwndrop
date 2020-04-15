var appLogin = Vue.component("app-login", {
    template: `
	<div class="row align-items-center">
		<div class="col-sm-1 col-md-2 col-lg-3"></div>
		<div class="col-sm-10 col-md-8 col-lg-6 v-middle">
			<div class="login-panel">
				<form>
					<div class="form-group row">
						<label for="login-username" class="col-sm-12 col-form-label">Username:</label>
						<div class="col-sm-12">
							<input
								type="text"
								class="form-control"
								id="login-username"
                                spellcheck="false"
                                autocomplete="off"
								v-model="Username"
                                name="username"
                                v-validate="'required'"
                                v-bind:class="{'form-control': true, 'error': errors.has('username') }"
								required
							>
                            <div v-show="errors.has('username')" class="form-error">{{ errors.first('username') }}</div>
						</div>
					</div>
					<div class="form-group row">
						<label for="login-password" class="col-sm-12 col-form-label">Password:</label>
						<div class="col-sm-12">
							<input
								type="password"
								class="form-control"
								id="login-password"
								spellcheck="false"
								v-model="Password"
                                name="password"
                                v-validate="'required'"
                                v-bind:class="{'form-control': true, 'error': errors.has('password') }"
								required
							>
                            <div v-show="errors.has('password')" class="form-error">{{ errors.first('password') }}</div>
						</div>
					</div>
					<div v-show="status != ''" class="form-group row">
						<div class="col-sm-12">
                            <div class="login-status">
                                {{ status }}
                            </div>
                        </div>
					</div>
					<hr>
					<div class="form-group row">
						<div class="col-sm-12">
							<button
								class="btn btn-primary btn-lg btn-login"
								type="submit"
                                :disabled="errors.any() || !isComplete"
								@click.prevent="login()"
							>Login</button>
						</div>
					</div>
				</form>
			</div>
		</div>
	</div>
    `,
    $_veeValidate: {
        validator: "new"
    },
	data() {
		return {
			url: Config.Hostname + Config.AdminDir + "/" + Config.ApiPath,
			Username: "",
            Password: "",
            status: "",
		};
	},
    computed: {
        isComplete () {
            return this.Username && this.Password;
        }
    },
	methods: {
		login() {
			if (this.Username == "" || this.Password == "") {
				return;
			}

			axios
				.post(
					this.url + "/login",
					{
						username: this.Username,
						password: this.Password
					},
					{
						headers: {
							"content-type": "application/json"
						}
					}
				)
				.then(response => {
                    console.log(response);
                    console.log(response.data.data.username);
					this.mainBus.$emit("loggedIn", response.data.data.username);
				})
				.catch(error => {
                    if (error.response.status == 401)
                        this.status = "Incorrent username or password"
                    else
                        this.status = "Internal server error"
					console.log(error);
				});
		}
	}
})