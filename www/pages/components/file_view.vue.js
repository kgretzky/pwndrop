var appFileView = Vue.component("app-file-view", {
    template: `
	<div>
		<b-modal
			v-model="editShow"
			id="edit-modal"
			title="Edit"
			size="lg"
			hide-header
			ok-title="Save"
			@ok.prevent="updateFile()"
		>
			<form>
				<div class="form-group row">
					<label for="edit-name" class="col-sm-3 col-form-label label-help">Name:
                        <i class="fas fa-question-circle label-qmark" v-tooltip:bottom="'Friendly name for your eyes only'"></i>
                    </label>
					<div class="col-sm-9">
						<input
							type="text"
							class="form-control"
							id="edit-name"
							spellcheck="false"
                            v-model="file_edit.name"
                            name="display-name"
                            v-validate="'required'"
                            v-bind:class="{'form-control': true, 'error': errors.has('display-name') }"
                        >
                        <div v-show="errors.has('display-name')" class="form-error"">{{ errors.first('display-name') }}</div>
					</div>
				</div>
				<div class="form-group row">
                    <label for="edit-mime" class="col-sm-3 col-form-label label-help">
                        <a class="help-link" href="https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types/Complete_list_of_MIME_types" target="_blank">
                        MIME Type:</a>
                        <i class="fas fa-question-circle label-qmark" v-tooltip:bottom="'File will be retrieved with the following MIME type'"></i>
                    </label>
					<div class="col-sm-9">
						<div class="input-group">
							<input
								type="text"
								ref="editMime"
								class="form-control"
								id="edit-mime"
								spellcheck="false"
                                v-model="file_edit.mime_type"
                                name="mime-type"
                                v-validate="'required'"
                                v-bind:class="{'form-control': true, 'error': errors.has('mime-type') }"
							>
							<div class="input-group-append">
								<button
									class="btn btn-secondary"
									type="button"
									@click="file_edit.mime_type = file_edit.orig_mime_type"
								>
									<i class="fas fa-undo"></i>
								</button>
                            </div>
                        </div>
                        <div v-show="errors.has('mime-type')" class="form-error">{{ errors.first('mime-type') }}</div>
					</div>
				</div>
				<div class="form-group row">
					<label for="edit-http-path" class="col-sm-3 col-form-label label-help">Path:
                        <i class="fas fa-question-circle label-qmark" v-tooltip:bottom="'URL path when sharing it over HTTP or WebDAV. Paths for WebDAV must be under a subdirectory (e.g. /subdir/payload.docx)'"></i>
                    </label>
					<div class="col-sm-9">
						<input
							type="text"
							class="form-control"
							id="edit-http-path"
							spellcheck="false"
                            v-model="file_edit.url_path"
                            name="url-path"
                            v-validate="'required'"
                            v-bind:class="{'form-control': true, 'error': errors.has('url-path') }"
                        >
                        <div v-show="errors.has('url-path')" class="form-error">{{ errors.first('url-path') }}</div>
					</div>
				</div>
				<div class="form-group row">
					<label for="edit-redirect-path" class="col-sm-3 col-form-label label-help">Redirect Path:
                        <i class="fas fa-question-circle label-qmark" v-tooltip:bottom="'URL path which the request will be redirected to. Useful if you want to spoof the extension (e.g. /subdir/payload.docx.exe)'"></i>
                    </label>
					<div class="col-sm-9">
                        <div class="input-group">
                            <input
                                type="text"
                                class="form-control"
                                id="edit-redirect-path"
                                spellcheck="false"
                                v-model="file_edit.redirect_path"
                                name="redirect-path"
                            >
                            <div class="input-group-append">
                                <button
                                    class="btn btn-secondary"
                                    type="button"
                                    @click="file_edit.redirect_path = file_edit.url_path"
                                >
                                    <i class="fas fa-copy"></i>
                                </button>
                            </div>
                        </div>
					</div>
				</div>
				<hr>
				<transition name="sub-modal-anim" mode="out-in">
					<div class="row" v-if="file_edit.sub_progress < 100" key="uploading">
						<div class="file-progress col">
							<div class="progress">
								<div
									class="progress-bar progress-bar-striped progress-bar-animated bg-success"
									role="progressbar"
									:style="{width: file_edit.sub_progress+'%'}"
									aria-valuemin="0"
									aria-valuemax="100"
									:aria-valuenow="file_edit.sub_progress"
								></div>
							</div>
						</div>
					</div>
					<div class="row" v-else-if="file_edit.ref_sub_file == 0" key="empty">
						<div class="sub-info">
							<small>Upload a facade file, which will be served instead of the original one, only when facade is enabled.</small>
						</div>
						<div id="sub-dropzone" :class="[isSubDragging ? 'drag' : '']">
							<span class="icon">
								<i class="fas fa-upload"></i>
							</span>
							<input
								type="file"
								@change="handleSubFile($event)"
								@dragover.prevent
								@dragenter="isSubDragging = true"
								@drop.prevent="handleSubDrop($event)"
								@dragleave="isSubDragging = false"
							>
						</div>
					</div>
                    <div class="sub-item" v-else key="uploaded">
                        <div class="sub-info">
                            <small>Facade file</small>
                        </div>
                        <div class="form-group row desc">
                            <label for="edit-sub-name" class="col-sm-3 col-form-label label-help">
                                Name:</a>
                                <i class="fas fa-question-circle label-qmark" v-tooltip:bottom="'Facade name for display purposes only'"></i>
                            </label>
							<div class="col">
                                <input type="text" class="form-control" spellcheck="false" v-model="file_edit.sub_name"
                                name="sub-name"
                                v-validate="'required'"
                                v-bind:class="{'form-control': true, 'error': errors.has('sub-name') }">
                                <div v-show="errors.has('sub-name')" class="form-error">{{ errors.first('sub-name') }}</div>
							</div>
							<div class="d-none d-sm-block col-auto shrink">
								<span class="fsize">{{ file_edit.sub_size | prettyBytes }}</span>
							</div>
							<div class="controls">
								<button
									class="btn btn-sm btn-danger btn-circle-sm"
									@click.prevent="deleteSubFile(file_edit.id, file_edit.ref_sub_file)"
								>
									<i class="fas fa-times"></i>
								</button>
							</div>
                        </div>
                        <div class="form-group row desc">
                            <label for="edit-sub-mime" class="col-sm-3 col-form-label label-help">
                                <a class="help-link" href="https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types/Complete_list_of_MIME_types" target="_blank">
                                MIME Type:</a>
                                <i class="fas fa-question-circle label-qmark" v-tooltip:bottom="'Facade file will be retrieved with the following MIME type'"></i>
                            </label>
                            <div class="col-sm-9">
                                <div class="input-group">
                                    <input
                                        type="text"
                                        ref="editSubMime"
                                        class="form-control"
                                        id="edit-sub-mime"
                                        spellcheck="false"
                                        v-model="file_edit.sub_mime_type"
                                        name="sub-mime-type"
                                        v-validate="'required'"
                                        v-bind:class="{'form-control': true, 'error': errors.has('sub-mime-type') }"
                                    >
                                    <div class="input-group-append">
                                        <button
                                            class="btn btn-secondary"
                                            type="button"
                                            @click="file_edit.sub_mime_type = file_edit.orig_mime_type"
                                        >
                                            <i class="fas fa-undo"></i>
                                        </button>
                                    </div>
                                </div>
                                <div v-show="errors.has('sub-mime-type')" class="form-error">{{ errors.first('sub-mime-type') }}</div>
                            </div>
                        </div>
					</div>
				</transition>
			</form>
		</b-modal>
		<div
			id="dropzone"
			:class="[isDragging ? 'drag' : '']"
			@dragover.prevent
			@dragenter="isDragging = true"
			@drop.prevent="handleDrop($event)"
			@dragleave="isDragging = false"
		></div>
		<div class="row row-file">
			<div class="col-xs-12 col-sm-8 offset-sm-2 col-md-6 offset-md-3">
				<form enctype="multipart/form-data" novalidate>
					<button class="btn btn-lg btn-primary btn-file btn-block">
						Upload
						<input type="file" multiple @change="handleFiles($event)">
					</button>
				</form>
			</div>
        </div>
        <div class="row row-info">
            <div class="server-status col-xs-12 col-sm-8 offset-sm-2 col-md-6 offset-md-3">
                free: <strong>{{ server_info.disk_free | prettyBytes }}</strong> &bull; used: <strong>{{ server_info.disk_used | prettyBytes }}</strong>
            </div>
        </div>
		<!-- 		<button @click="doShuffle()">Shuffle</button>
		-->
		<transition-group name="upload-list">
			<div class="row upload-block" v-for="upload in uploads" :key="upload.key">
				<app-file
					:file="upload"
					@editFile="editFile"
					@deleteFile="deleteFile"
					@enableFile="enableFile"
					@pauseFile="pauseFile"
				></app-file>
			</div>
		</transition-group>
    </div>
    `,
    $_veeValidate: {
        validator: "new"
    },
	data() {
		return {
			url: Config.Hostname + Config.AdminDir + "/" + Config.ApiPath,
			isDragging: false,
			isSubDragging: false,
			editShow: false,
			uploads: [],
            next_key: 0,
			file_edit: {
				create_time: 0,
				fsize: 0,
				id: 0,
                mime_type: "",
                sub_mime_type: "",
				name: "",
				orig_mime_type: "",
				ref_sub_file: 0,
				sub_ctime: 0,
				sub_name: "",
				sub_progress: 100,
				sub_size: 0,
                url_path: "",
                redirect_path: "",
				wdav_path: ""
            },
            server_info: {
                disk_free: 0,
                disk_used: 0
            }
		};
    },
    computed: {
        isComplete () {
            return this.file_edit.name && this.file_edit.mime_type && this.file_edit.url_path && (this.file_edit.ref_sub_file == 0 || (this.file_edit.sub_name && this.file_edit.sub_mime_type));
        }
    },
	methods: {
		doShuffle() {
			this.uploads = _.shuffle(this.uploads);
		},
		handleFiles($event) {
			console.log($event.target.files);
			var files = [];
			for (var i = 0; i < $event.target.files.length; i++) {
				var file = $event.target.files[i];
				files.push(file);
			}
			if (files.length > 0) {
				this.uploadFiles(files);
			}
		},
		handleDrop($event) {
			this.isDragging = false;
			var files = [];
			if ($event.dataTransfer.items) {
				for (var i = 0; i < $event.dataTransfer.items.length; i++) {
					var file = $event.dataTransfer.items[i].getAsFile();
					files.push(file);
				}
			} else {
				for (var i = 0; i < $event.dataTransfer.items.length; i++) {
					var file = $event.dataTransfer.files[i];
					files.push(file);
				}
			}
			if (files.length > 0) {
				this.uploadFiles(files);
			}
		},
		uploadFiles(files) {
			var vm = this;
			console.log(files);
			var ctime = new Date().getTime();
			for (var i = 0; i < files.length; i++) {
				const file = files[i];
				const item_id = ctime + i;
				var item = {
					id: item_id,
					name: file.name,
					fsize: file.size,
                    mime_type: file.type,
                    sub_mime_type: file.sub_mime_type,
					orig_mime_type: file.type,
                    url_path: "",
                    redirect_path: "",
					wdav_path: "",
					progress: 0,
					key: this.next_key,
					is_enabled: true,
					is_paused: false,
                    sub_name: "",
                    sub_file: null,
				};
				vm.uploads.push(item);
				this.next_key += 1;

				console.log("uploading: " + file.name);
				var formData = new FormData();
				formData.append("file", file);
				axios
					.post(vm.url + "/files", formData, {
						headers: {
							"content-type": "multipart/form-data"
						},
						onUploadProgress(progressEvent) {
							var i = vm.findFileIndexById(item_id);
							if (i != -1) {
								vm.uploads[i].progress = Math.floor(
									(progressEvent.loaded /
										progressEvent.total) *
										100
								);
							}
						}
					})
					.then(response => {
						console.log(response);
						console.log("item_id: " + item_id);
						var i = vm.findFileIndexById(item_id);
						if (i != -1) {
							//vm.uploads[i] = response.data.data;
							var it = response.data.data;
							var ut = vm.uploads[i];

							ut.id = it.id;
                            ut.url_path = it.url_path;
                            ut.redirect_path = it.redirect_path;
                            ut.wdav_path = it.wdav_path;
                            ut.mime_type = it.mime_type;
                            ut.sub_mime_type = it.sub_mime_type;
                            ut.orig_mime_type = it.orig_mime_type;
							ut.progress = 100;

							//it.progress = 100;
							//it.key = vm.uploads[i].key;
							//console.log("it.key:", it.key);
							//vm.uploads[i] = it;
							//vm.uploads[i].key = vm.next_key;
							//vm.next_key += 1;
							//vm.uploads[i].id = item_id;
							//vm.uploads.splice(i, 1, it);
                            //vm.uploads[i].progress = 100;
                            
							//vm.uploads.sort((a, b) => a.key - b.key);
                        }
                        this.syncServerInfo();

						console.log(vm.uploads);
					})
					.catch(error => {
                        console.log(error)

                        var i = vm.findFileIndexById(item_id);
						if (i != -1) {
						    vm.uploads.splice(i, 1);
                        }
                    });
			}
		},
		handleSubFile($event) {
			console.log($event.target.files);
			var files = [];
			for (var i = 0; i < $event.target.files.length; i++) {
				var file = $event.target.files[i];
				files.push(file);
				break;
			}
			if (files.length > 0) {
				this.uploadSubFiles(this.file_edit.id, files);
			}
		},
		handleSubDrop($event) {
			this.isSubDragging = false;
			var files = [];
			if ($event.dataTransfer.items) {
				for (var i = 0; i < $event.dataTransfer.items.length; i++) {
					var file = $event.dataTransfer.items[i].getAsFile();
					files.push(file);
					break;
				}
			} else {
				for (var i = 0; i < $event.dataTransfer.items.length; i++) {
					var file = $event.dataTransfer.files[i];
					files.push(file);
					break;
				}
			}
			if (files.length > 0) {
				this.uploadSubFiles(this.file_edit.id, files);
			}
		},
		uploadSubFiles(parent_id, files) {
			var vm = this;
			console.log(files);
			var ctime = new Date().getTime();
			for (var i = 0; i < files.length; i++) {
				const file = files[i];

				console.log("uploading sub: " + file.name);
				var formData = new FormData();
				formData.append("file", file);
				axios
					.post(vm.url + "/files/" + parent_id + "/sub", formData, {
						headers: {
							"content-type": "multipart/form-data"
						},
						onUploadProgress(progressEvent) {
							console.log(progressEvent);
							vm.file_edit.sub_progress = Math.floor(
								(progressEvent.loaded / progressEvent.total) *
									100
							);
						}
					})
					.then(response => {
						console.log(response);
						vm.file_edit.sub_progress = 100;
						var it = response.data.data;
						vm.file_edit.ref_sub_file = it.id;
						vm.file_edit.sub_name = it.name;
						vm.file_edit.sub_size = it.fsize;
						vm.file_edit.sub_ctime = it.create_time;

						var i = this.findFileIndexById(parent_id);
						if (i != -1) {
							var f = this.uploads[i];
							f.ref_sub_file = it.id;
							f.sub_file = {
								create_time: it.create_time,
								fid: parent_id,
								fname: it.fname,
								fsize: it.fsize,
								id: it.id,
								name: it.name,
								uid: it.uid
                            };
                            f.sub_name = it.name;
                            console.log(f);
                        }
                        
                        this.syncServerInfo();
					})
					.catch(error => console.log(error));
			}
		},
		editFile(id) {
			var i = this.findFileIndexById(id);
			if (i == -1) {
				console.log("file not found: " + id);
				return;
			}
			this.file_edit.id = id;
			this.file_edit.name = this.uploads[i].name;
            this.file_edit.mime_type = this.uploads[i].mime_type;
            this.file_edit.sub_mime_type = this.uploads[i].sub_mime_type;
			this.file_edit.orig_mime_type = this.uploads[i].orig_mime_type;
			this.file_edit.url_path = this.uploads[i].url_path;
			this.file_edit.redirect_path = this.uploads[i].redirect_path;
			this.file_edit.wdav_path = this.uploads[i].wdav_path;
			this.file_edit.ref_sub_file = this.uploads[i].ref_sub_file;
			this.file_edit.sub_name = "<unknown>";
			this.file_edit.sub_size = 0;
			this.file_edit.sub_ctime = 0;
			this.file_edit.sub_progress = 100;
			if (this.uploads[i].sub_file) {
				this.file_edit.sub_name = this.uploads[i].sub_name;
				this.file_edit.sub_size = this.uploads[i].sub_file.fsize;
				this.file_edit.sub_ctime = this.uploads[i].sub_file.create_time;
			} else {
				this.file_edit.ref_sub_file = 0;
			}
			//this.$refs.editModal.modal("show");
			//this.edit_show = true;
			this.$bvModal.show("edit-modal");
		},
		updateFile() {
            var vm = this;
			if (!this.file_edit) {
				return;
            }
            if (!this.isComplete) {
                return;
            }
			var id = this.file_edit.id;
			axios
				.put(
					this.url + "/files/" + id,
					{
						name: this.file_edit.name,
						url_path: this.file_edit.url_path,
						redirect_path: this.file_edit.redirect_path,
                        mime_type: this.file_edit.mime_type,
                        sub_mime_type: this.file_edit.sub_mime_type,
						sub_name: this.file_edit.sub_name
					},
					{
						headers: {
							"content-type": "application/json"
						}
					}
				)
				.then(response => {
					console.log(response);
					vm.$bvModal.hide("edit-modal");

					var i = this.findFileIndexById(id);
					if (i != -1) {
						var f = response.data.data;
						vm.uploads[i].name = f.name;
						vm.uploads[i].sub_name = f.sub_name;
						vm.uploads[i].url_path = f.url_path;
						vm.uploads[i].redirect_path = f.redirect_path;
                        vm.uploads[i].mime_type = f.mime_type;
                        vm.uploads[i].sub_mime_type = f.sub_mime_type;                       
					}
				})
				.catch(error => {
					console.log(error);
				});
		},
		deleteFile(id) {
			axios
				.delete(this.url + "/files/" + id)
				.then(response => {
					console.log(response);
					var i = this.findFileIndexById(id);
					if (i != -1) {
						this.uploads.splice(i, 1);
                    }
                    this.syncServerInfo();
				})
				.catch(error => {
					console.log(error);
				});
		},
		enableFile(id) {
			var i = this.findFileIndexById(id);
			if (i == -1) {
				return;
			}
			var api = "/enable";
			if (this.uploads[i].is_enabled) {
				api = "/disable";
			}

			axios
				.get(this.url + "/files/" + id + api)
				.then(response => {
					console.log(response);
					var i = this.findFileIndexById(id);
					if (i != -1) {
						var f = response.data.data;
						this.uploads[i].is_enabled = f.is_enabled;
						this.uploads[i].is_paused = f.is_paused;
					}
				})
				.catch(error => {
					console.log(error);
				});
		},
		pauseFile(id) {
			var i = this.findFileIndexById(id);
			if (i == -1) {
				return;
			}
			if (this.uploads[i].sub_file == null) {
				this.uploads[i].is_paused = false;
				return;
			}

			var api = "/pause";
			if (this.uploads[i].is_paused) {
				api = "/unpause";
			}

			axios
				.get(this.url + "/files/" + id + api)
				.then(response => {
					console.log(response);
					var i = this.findFileIndexById(id);
					if (i != -1) {
						var f = response.data.data;
						this.uploads[i].is_enabled = f.is_enabled;
						this.uploads[i].is_paused = f.is_paused;
					}
				})
				.catch(error => {
					console.log(error);
				});
		},
		deleteSubFile(parent_id, sub_id) {
			axios
				.delete(this.url + "/files/" + parent_id + "/sub/" + sub_id)
				.then(response => {
					console.log(response);

					this.file_edit.ref_sub_file = 0;
					var i = this.findFileIndexById(parent_id);
					if (i != -1) {
						this.uploads[i].ref_sub_file = 0;
						this.uploads[i].sub_name = "";
						this.uploads[i].sub_file = null;
						this.uploads[i].is_paused = false;
                    }
                    this.syncServerInfo();
				})
				.catch(error => {
					console.log(error);
				});
		},
		findFileIndexById(id) {
			var ret = -1;
			this.uploads.forEach(function(it, i) {
				if (it.id == id) {
					ret = i;
				}
			});
			return ret;
		},
		refresh() {
			var t = this;
			axios
				.get(t.url + "/files")
				.then(response => {
					console.log(response);
					var files = response.data.data.uploads;
					t.uploads = [];
					for (var i = 0; i < files.length; i++) {
						t.uploads.push(files[i]);
						t.uploads[i].key = i;
					}
					t.next_key = i + 1;
				})
				.catch(error => console.log(error));
		},
        syncServerInfo() {
            axios
                .get(this.url + "/server_info")
                .then(response => {
                    var r = response.data.data;

                    this.server_info.disk_free = r.disk_free;
                    this.server_info.disk_used = r.disk_used;

                    console.log(r);
                })
                .catch(error => {
                    console.log(error);
                });
        }
	},
	created() {
		var t = this;
		window.addEventListener("dragenter", function(e) {
			if (!t.editShow) {
				t.isDragging = true;
			}
        });
        this.syncServerInfo();
        this.refresh();
	}
})