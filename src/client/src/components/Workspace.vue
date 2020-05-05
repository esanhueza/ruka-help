<template>
  <div id="workspace-wrapper" class="">
    <div id="workspace-toolbar" class="">
      <unicon v-on:click="toolbarSelection = 'chat'" class="toolbar-btn" name="comments-alt" v-bind:fill="toolbarSelection == 'chat' ? 'royalblue' : 'gray'" />
      <unicon v-on:click="toolbarSelection = 'users'" class="toolbar-btn" name="users-alt" v-bind:fill="toolbarSelection == 'users' ? 'royalblue' : 'gray'" />
    </div>
    <div id="workspace">
      <div class="secondary-tool-wrapper">
        <UserList v-if="this.toolbarSelection == 'users'" v-bind:users="this.users"/>
        <Chat v-if="this.toolbarSelection == 'chat'" v-on:message-sent="this.sendChatMessage" v-bind:status="this.status" v-bind:user="this.user" v-bind:history="this.history"/>
      </div>
      <div class="board-wrapper">
        <Board v-bind:code="this.code"/>
      </div>
    </div>
  </div>
</template>

<script>
import Board from '../components/Board.vue'
import UserList from '../components/tools/UserList.vue'
import Chat from '../components/chat/Chat.vue'

export default {
  name: 'Workspace',
  components: {
      Board,
      Chat,
      UserList,
  },
  data: function(){
    return {
      chatStatus: 'disconnected',
      toolbarSelection: 'chat',
      status: "disconnected",
      socket: null,
      user: null,
      history: [],
      users: [],
      code: "",
      session: null,
    };
  },
  beforeDestroy(){
    this.disconnect();
  },
  created() {
    window.addEventListener('unload', this.unload);
    if (!sessionStorage.getItem("UserID") ){
      this.$router.push({name: 'home'});
    }
    else{
      this.session = {
        id: sessionStorage.getItem("SessionID")
      };
      this.user = { 
          id: sessionStorage.getItem("UserID"),
          username: sessionStorage.getItem("Username"),
          status: 'online'
      };
      this.users.push(this.user);
      this.connect();
    }
  },
  methods: {
    unload(){
      console.log("unload");
      this.disconnect();
    },
    connect() {
      const that = this;
      this.socket = new WebSocket("ws://localhost:9999/ws?token=" + this.user.id);
      this.socket.onopen = () => {
        this.status = "connected";
        console.log("on socket opened");

        // start chat
        this.socket.send(JSON.stringify({
          UserID: this.user.id,
          SessionID: this.session.id,
          Type: "greetings"
        }));
        this.chatStatus = true;
        
        // handle incoming messages
        this.socket.onmessage = ({ data }) => {
          let message = JSON.parse(data);
          switch (message.Type){
            case "message":
              that.handleChatMessage(message);
              break;
            case "system":
              that.handleSystemNotification(message);
              break;
            case "sync":
              that.handleSyncMessage(message);
              break;
            default:
              console.log("Handler missing for messages of type: ", message.Type);
          }
        };
      };
    },
    disconnect () {
      this.socket.close()
      this.status = 'disconnected'
      console.log('WebSocket disconnected')
    },
    sendChatMessage(data){
      if (this.chatStatus === 'disconnected') return;
      this.socket.send(JSON.stringify(Object.assign(data, {
        SessionID: this.session.id,
        UserID: this.user.id,
      })));
    },
    handleSystemNotification(data){
      console.log("handleSystemNotification"),
      this.history.push({
          id: data.ID,
          type: data.Type,
          messageClass: data.Type,
          content: data.Content,
          user: {
            id: "",
            username: ""
          }
      });
    },
    handleChatMessage(data){
      console.log("handleChatMessage"),
      this.history.push({
          id: data.ID,
          user: {
            id: data.User.ID,
            username: data.User.Username
          },
          messageClass: data.User.ID == this.user.id ? "own" : "other",
          type: data.Type,
          content: data.Content,
      });
    },
    handleSyncMessage(data){
      let that = this;
      this.history = [];
      data.Chat.map(function(d){
        that.history.push({
          id: d.ID,
          user: {
            id: d.User.ID,
            username: d.User.Username
          },
          messageClass: d.Type === "system" ? d.Type : (d.User.ID == that.user.id ? "own" : "other"),
          type: d.Type,
          content: d.Content,
        });
      });
      this.users = [];
      if (data.Users){
        data.Users.map(function(d){
          that.users.push({
            id: d.ID,
            username: d.Username,
            status: "online",
          });
        });
      }
    },
    enableChat(){
      this.socket.send(JSON.stringify({
          UserID: this.user.id,
          SessionID: this.session.id,
          Type: "greetings"
        }));
        this.chatStatus = 'connected';
    },
    disableChat(){
      this.socket.send(JSON.stringify({
          UserID: this.user.id,
          SessionID: this.session.id,
          Type: "goodbye"
        }));
        this.chatStatus = 'disconnected';
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
#workspace-wrapper{
    height: 100%;
    margin: 0px;
    width: 100%;
    display: flex;
    flex-direction: row;
    justify-content: center;
    align-items: start;
}
#workspace{
    display: flex;
    flex-direction: row;
    width: 100%;
    height: 100%;
}
.secondary-tool-wrapper{
    width: 30%;
    border: 1px solid black;
    height: 80%;
}
.board-wrapper{
    margin-left: 20px;
    width: 100%;
    border: 1px solid black;
    height: 80%;
}
#workspace-toolbar{
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-items: start;
  background: #d9d9d9;
  border: 1px solid black;
  border-right: unset;
}
#workspace-toolbar .toolbar-btn{
  margin: 10px 15px 10px 15px;
}

#workspace-toolbar .toolbar-btn:after{
  position: absolute;
  left: 54px;
  width: 4px;
  background: #d9d9d9;

}
#workspace-toolbar .toolbar-btn.selected{
}

</style>
