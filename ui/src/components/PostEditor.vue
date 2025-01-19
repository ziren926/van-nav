<template>
  <div class="post-editor">
    <a-form :model="postForm" layout="vertical">
      <a-form-item label="标题">
        <a-input v-model:value="postForm.post_title" placeholder="请输入标题" />
      </a-form-item>
      <a-form-item label="内容">
        <a-textarea
          v-model:value="postForm.post_content"
          :rows="6"
          placeholder="请输入内容"
        />
      </a-form-item>
      <a-form-item>
        <a-button type="primary" @click="handleSubmit">保存</a-button>
      </a-form-item>
    </a-form>
  </div>
</template>

<script>
import { defineComponent, ref, onMounted } from 'vue';
import { message } from 'ant-design-vue';
import { useRouter, useRoute } from 'vue-router';

export default defineComponent({
  name: 'PostEditor',
  setup() {
    const router = useRouter();
    const route = useRoute();
    const postForm = ref({
      post_title: '',
      post_content: ''
    });

    const fetchPost = async () => {
      try {
        const response = await fetch(`/api/admin/tool/${route.params.id}/post`, {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        });
        if (response.ok) {
          const data = await response.json();
          postForm.value = data;
        }
      } catch (error) {
        console.error('获取帖子失败:', error);
      }
    };

    const handleSubmit = async () => {
      try {
        const response = await fetch(`/api/admin/tool/${route.params.id}/post`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          },
          body: JSON.stringify(postForm.value)
        });

        if (response.ok) {
          message.success('保存成功');
          router.push('/admin/tools');
        } else {
          message.error('保存失败');
        }
      } catch (error) {
        message.error('保存失败: ' + error.message);
      }
    };

    onMounted(() => {
      fetchPost();
    });

    return {
      postForm,
      handleSubmit
    };
  }
});
</script>

<style scoped>
.post-editor {
  padding: 24px;
  background: #fff;
  border-radius: 2px;
}
</style>