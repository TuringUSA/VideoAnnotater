<template>
  <b-collapse
    class="card my-3"
    animation="slide"
    :open="false"
    :aria-id="uuid"
  >
    <template #trigger="props">
      <div
        class="card-header"
        role="button"
        aria-controls="contentIdForA11y3"
        :aria-expanded="props.open">
        <p class="card-header-title">
          {{ fileName }}
        </p>
        <a class="card-header-icon">
          <b-icon
            :icon="props.open ? 'menu-down' : 'menu-up'">
          </b-icon>
        </a>
      </div>
    </template>

    <div class="card-content">
      <div class="is-flex columns">
        <div class="is-half column">
          <h3>Video</h3>
          <video controls style="max-height: 500px" ref="video">
            <source :src="link" type="video/mp4">
            Your browser does not support the video tag.
          </video>
        </div>
        <div class="column is-half has-text-left">
          <h3>Features</h3>
          <b-field grouped group-multiline>

            <template v-for="feature in features">
            <b-taglist attached v-for="timestamp in feature.Timestamps" class="mx-1 mb-0">
              <b-tag type="is-dark">
                {{ feature.Description }}
              </b-tag>
              <b-tag type="is-info" @click="moveToTime(timestamp.Start)" style="cursor: pointer">{{timestamp.Start}}</b-tag>
              <b-tag type="end-time-tag" @click="moveToTime(timestamp.End)"  style="cursor: pointer">{{msToTime(timestamp.End)}}</b-tag>
            </b-taglist>
            </template>
          </b-field>
        </div>

      </div>
    </div>

    <footer class="card-footer">
      <a class="card-footer-item" @click="newTab">Open in new tab</a>
    </footer>
  </b-collapse>
</template>

<script>
export default {
  name: "VideoResult",
  props: {
    fileName: {
      type: String,
      default: ""
    },
    uuid: {
      type: String,
      required: true
    },
    link: {
      type: String,
      required: true
    },
    features: {
      type: Array
    }
  },
  methods: {
    newTab() {
      let linkEl = document.createElement("a");
      linkEl.download = name;
      linkEl.href = this.link;
      linkEl.target = "_blank";
      document.body.appendChild(linkEl);
      linkEl.click();
      document.body.removeChild(linkEl);
    },
    msToTime(duration) {
      let milliseconds = Math.floor((duration % 1000) / 100),
        seconds = Math.floor((duration / 1000) % 60),
        minutes = Math.floor((duration / (1000 * 60)) % 60),
        hours = Math.floor((duration / (1000 * 60 * 60)) % 24);

      hours = (hours < 10) ? "0" + hours : hours;
      minutes = (minutes < 10) ? "0" + minutes : minutes;
      seconds = (seconds < 10) ? "0" + seconds : seconds;

      return hours + ":" + minutes + ":" + seconds + "." + Math.floor(milliseconds);
    },

    moveToTime(time){
      this.$refs.video.currentTime = time;
    }

  }
}
</script>

<style scoped lang="scss">
  .end-time-tag{
    background-color: #136bcc;
    color: white;
  }
</style>
