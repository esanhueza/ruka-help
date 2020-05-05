<template>
    <div id="chat">
      <div id="messages" ref="messageHistory">
        <Message
          v-for="message in history"
          v-bind:key="message.id"
          v-bind:content="message.content"
          v-bind:type="message.type"
          v-bind:messageClass="message.messageClass"
          v-bind:username="message.user.username"
        ></Message>
      </div>
      <hr />
      <div id="new-message">
        <textarea
          v-model="currentMessage"
          id="textarea-message"
          name=""
          cols="30"
          rows="2"
        ></textarea>
        <button v-on:click="sendMessage" type="submit" id="btn-send-message">
          Enviar
        </button>
      </div>
    </div>
</template>

<script>
import Message from "./Message.vue";

export default {
  name: "Chat",
  props: [
    'history',
    'status',
    'user',
  ],
  data: function () {
    return {
      currentMessage: "",
    };
  },
  components: {
    Message,
  },
  watch: {
      history: function(){
        this.$nextTick(function() {
            var container = this.$refs.messageHistory;
            container.scrollTop = container.scrollHeight + 120;
        });
      }
  },
  methods: {
    sendMessage: function (event) {
      event.preventDefault();
      if (this.currentMessage == ""){
          return;
      }
      this.$emit('message-sent', {
        "Type": "message",
        "Content": this.currentMessage
      });
      this.currentMessage = "";
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
#chat-wrapper {
  /* background-color: blue; */
}
#chat {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #d9d9d9;
}
.message {
  margin: 20px 15px 0 15px;
  /* 
    border: 1px solid black;
     */
  width: calc(100% - 90px);
  padding: 10px 0 0 10px;
}
#messages {
  display: flex;
  flex-direction: column;
  flex-grow: 20;
  overflow-y: auto;
}
#new-message {
  flex-grow: 1;
}
#textarea-message {
  width: 90%;
}
#btn-send-message {
  width: 91%;
  margin: 0px;
  padding: 0px;
}
</style>
