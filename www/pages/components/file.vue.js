var appFile = Vue.component("app-file", {
    template: `
	<div ref="uploadItem" class="upload-item" :class="{'upload-disabled': !file.is_enabled, 'upload-paused': file.is_paused}">
		<div class="row upload-desc">
			<div class="settings">
                <span class="btn-col">
				<button class="btn btn-sm btn-primary btn-circle-sm" @click="$emit('editFile', file.id)" v-tooltip:top="'Change file settings'">
					<i class="fas fa-cog"></i>
				</button>
                </span>
                <span class="btn-col">
				<button class="btn btn-sm btn-circle-sm" :class="{'btn-outline-secondary': !file.is_enabled, 'btn-success': file.is_enabled}" @click="$emit('enableFile', file.id)" v-tooltip:top="'Make file available for download'">
					<i class="fas fa-power-off"></i>
				</button>
                </span>
                <span class="btn-col">
				<button class="btn btn-sm btn-circle-sm" :class="{'btn-outline-secondary': !file.is_paused, 'btn-gray': file.is_paused}" @click="$emit('pauseFile', file.id)" v-tooltip:top="'Enable the facade and serve the facade file instead of the original one'">
					<i class="fas fa-mask"></i>
				</button>
                </span>
			</div>
			<div class="col clip trans" :class="{'text-dim': file.is_paused}">
				<span class="title">{{ file.name }}</span>
			</div>
            <div v-if="file.sub_file != null && !file.is_paused" class="col-auto shrink">
                <i class="fas fa-arrow-left"></i>
            </div>
            <div v-else-if="file.sub_file != null" class="col-auto shrink">
                <i class="fas fa-arrow-right"></i>
            </div>
			<div v-if="file.sub_file != null" class="col clip trans" :class="{'text-dim': !file.is_paused}">
				<span class="title">{{ file.sub_name }}</span>
			</div>
			<div class="d-none d-sm-block col-auto shrink text-right clip">
				<span class="fsize">{{ file.fsize | prettyBytes }}</span>
			</div>
			<div class="controls">
				<button class="btn btn-sm btn-danger btn-circle-sm" @click="deleteItem(file.id)">
					<i class="fas fa-times"></i>
				</button>
			</div>
		</div>
		<div class="row upload-info" v-show="!file.progress || file.progress == 100">
			<div class="col-auto shrink clip">
				<span class="btn-col">
					<a class="btn-copy" ref="copyUrl" href @click.prevent="copyHttpUrl()">
						<button class="btn btn-sm btn-outline-success btn-copy-link" v-tooltip:bottom="'Copy HTTP link to clipboard'">
							<i class="fas fa-copy" style="margin-right: 5px"></i>HTTP
						</button>
					</a>
				</span>
				<span class="btn-col">
					<a class="btn-copy" ref="copyWebdavUrl" href @click.prevent="copyWebdavUrl()">
						<button class="btn btn-sm btn-outline-success btn-copy-link" v-tooltip:bottom="'Copy WebDAV link to clipboard'">
							<i class="fas fa-copy" style="margin-right: 5px"></i>WebDAV
						</button>
					</a>
				</span>
			</div>
            <div class="col-auto grow trans" :class="{'text-lg': file.is_paused}">
                <small>{{ file.url_path }}</small>
            </div>
			<div class="d-none d-sm-block col text-right clip trans" :class="{'text-lg': file.is_paused}">
				<small>{{ file.mime_type }}</small>
			</div>
		</div>
		<div class="row">
			<div class="file-progress col" v-if="file.progress < 100">
				<div class="progress">
					<div
						class="progress-bar progress-bar-striped progress-bar-animated bg-success"
						role="progressbar"
						:style="{width: file.progress+'%'}"
						aria-valuemin="0"
						aria-valuemax="100"
						:aria-valuenow="file.progress"
					></div>
				</div>
			</div>
		</div>
	</div>
    `,
    $_veeValidate: {
        validator: "new"
    },
	props: ["file"],
	methods: {
		copyHttpUrl() {
			var l = window.location;
			var url = l.protocol + "//" + l.hostname;
			if (l.port != "" && (l.port != 443 && l.port != 80)) {
				url += ":" + l.port;
			}
			url += escape(this.file.url_path);
			this.$refs.copyUrl.setAttribute("data-clipboard-text", url);
		},
		copyWebdavUrl() {
			var l = window.location;
			var url = "\\\\" + l.hostname + "@80";
			/*if (l.port != 443 && l.port != 80) {
				url += "@" + l.port;
			}*/
			url += escape(this.file.url_path).replace(/\//g, "\\");
			this.$refs.copyWebdavUrl.setAttribute("data-clipboard-text", url);
		},
		deleteItem(id) {
			this.$refs.uploadItem.style.width =
				this.$refs.uploadItem.offsetWidth + "px";
			this.$emit("deleteFile", id);
		}
    }
})