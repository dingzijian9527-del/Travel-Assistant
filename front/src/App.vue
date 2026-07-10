<script>
export default {
  onLaunch() {
    console.log("Travel Assistant launched");
  }
};
</script>

<style>
page {
  min-height: 100%;
  background: #e5eeeb;
  color: #102033;
  font-family: -apple-system, BlinkMacSystemFont, "PingFang SC", "Helvetica Neue", Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
}

view,
text,
button,
input {
  box-sizing: border-box;
}

image {
  display: block;
}

button {
  padding: 0;
  margin: 0;
  line-height: 1;
  background: transparent;
  border-radius: 0;
  font-family: inherit;
}

button::after {
  border: 0;
}

.page {
  min-height: 100vh;
  background: #e5eeeb;
}

.phone-shell {
  width: 100%;
  max-width: 430px;
  margin: 0 auto;
  background: linear-gradient(180deg, #f8fbf9 0%, #eef5f2 46%, #f9faf7 100%);
}
</style>
