export default defineTask({
    name: "sync:models",
    run() {
        return performSync();
    },
});
