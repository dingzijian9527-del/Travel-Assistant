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
  background: #f8fafc;
  color: #0f172a;
  font-family: -apple-system, BlinkMacSystemFont, "PingFang SC", "Helvetica Neue", Arial, sans-serif;
}

view,
text,
button,
input {
  box-sizing: border-box;
}

button {
  padding: 0;
  margin: 0;
  line-height: 1;
  background: transparent;
}

button::after {
  border: 0;
}
</style>
