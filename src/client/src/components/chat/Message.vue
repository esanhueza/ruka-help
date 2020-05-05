<template>
    <div class="message" :class="messageClass">
        <div v-if="type === 'message'" class="message-owner">{{username}}:</div>
        <div class="message-content">
            {{ content }}
        </div>
    </div>
</template>

<script>
export default {
  name: 'Message',
  props: [
      'id',
      'content',
      'type',
      'messageClass',
      'username'
  ]
}
</script>

<style scoped>
.message{
    border-radius: 0.3em;
    font-size: 16px;
}
.message-owner{
    font-size: 14px;
    text-align: left;
    margin: 0.2em 0 0.2em 0;
}
.message-content{
    text-align: left;
    margin: 0.2em 0 0.4em 0;
}

.message.other{
    background: #ffffff;
    align-self: flex-start;
    text-align: left;
}
.message.own{
    background: #bababa;
    align-self: flex-end;
    text-align: right;
}

.message.system{
    background: #dbfbf5;
    align-self: center;
    text-align: center;
}

.message.system .message-content{
    text-align: center;
}
</style>
