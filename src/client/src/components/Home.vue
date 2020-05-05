<template>
    <div>
        <h1>Home</h1>
        <div>
            <input v-model="username" type="text" placeholder="Username">
        </div>
        <div>
            <input v-model="room" type="text" placeholder="Room">
        </div>
        <div>
            <button v-on:click="openRoom" type="button">Enter</button>
        </div>
    </div>
</template>

<script>


export default {
  name: 'Home',
  data: function(){
    return {
        username: "",
        room: "",
    }
  },
  methods: {
    openRoom(){
        console.log(this.$router);
        let that = this;
              // POST message user enetered to our server
        this.axios.post('http://localhost:9999/session/new', {Username: this.username, SessionID: this.room})
        .then(function(response){
                let data = response.data;
                sessionStorage.setItem("UserID", data.UserID);
                sessionStorage.setItem("SessionID", data.SessionID);
                sessionStorage.setItem("Username", data.Username);
                that.$router.push({name: 'workspace', params: { id: data.SessionID }});
        })
        .catch(function (error) {
            console.log(error)
        })
        
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
