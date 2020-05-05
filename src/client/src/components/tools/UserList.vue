<template>
    <div class="user-list" ref="userList">
    <UserListItem
        v-for="item in users"
        v-bind:key="item.id"
        v-bind:id="item.id"
        v-bind:username="item.username"
        v-bind:status="item.status"
    ></UserListItem>
    </div>
</template>

<script>
import UserListItem from "./UserListItem.vue";

export default {
  name: "UserList",
  props: [
    'users',
  ],
  data: function () {
    return {
    };
  },
  components: {
    UserListItem,
  },
  methods: {
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.user-list {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #d9d9d9;
}
</style>
