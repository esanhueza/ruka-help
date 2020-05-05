<template>
    <div class="user-list-item">
        <span>{{ username }}</span> <span class="user-status" v-bind:class="status">{{ status }}</span>
    </div>
</template>

<script>
export default {
  name: 'UserListItem',
  props: [
      'id',
      'username',
      'status'
  ]
}
</script>

<style scoped>
.user-list-item {
  padding: 5px 10px 5px 10px;
  height: 50px;
  align-items: middle;
  display: flex;
  flex-direction: row;
  border-bottom: 1px solid black;
  align-items: center;
  justify-content: space-between;
}

.user-list-item .user-status.online{
    color: green;
}

.user-list-item .user-status.offline{
    color: gray;
}

</style>
