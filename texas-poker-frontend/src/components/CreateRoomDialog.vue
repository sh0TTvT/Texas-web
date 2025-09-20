<!--
创建房间对话框组件
作用：提供创建新房间的表单界面
-->

<template>
  <q-card style="min-width: 400px; max-width: 500px">
    <q-card-section>
      <div class="text-h6">创建房间</div>
      <div class="text-subtitle2 text-grey-7">设置房间参数</div>
    </q-card-section>

    <q-card-section class="q-pt-none">
      <q-form @submit="handleSubmit" class="create-room-form">
        <!-- 房间名称 -->
        <q-input
          v-model="form.name"
          label="房间名称"
          :rules="nameRules"
          outlined
          maxlength="50"
          counter
          class="form-field"
        >
          <template v-slot:prepend>
            <q-icon name="title" />
          </template>
        </q-input>

        <!-- 筹码级别 -->
        <q-select
          v-model="form.chipLevel"
          label="筹码级别"
          :options="chipLevelOptions"
          outlined
          emit-value
          map-options
          class="form-field"
          @update:model-value="updateChipSettings"
        >
          <template v-slot:prepend>
            <q-icon name="star" />
          </template>
        </q-select>

        <!-- 最低筹码要求 -->
        <q-input
          v-model.number="form.minChips"
          label="最低筹码要求"
          type="number"
          :rules="minChipsRules"
          outlined
          class="form-field"
          :hint="`建议值: ${recommendedMinChips}`"
        >
          <template v-slot:prepend>
            <q-icon name="account_balance" />
          </template>
        </q-input>

        <!-- 盲注设置 -->
        <div class="blind-settings">
          <div class="blind-title">盲注设置</div>
          
          <div class="blind-inputs">
            <q-input
              v-model.number="form.smallBlind"
              label="小盲注"
              type="number"
              :rules="blindRules"
              outlined
              class="blind-input"
              @update:model-value="updateBigBlind"
            >
              <template v-slot:prepend>
                <q-icon name="remove_red_eye" />
              </template>
            </q-input>

            <q-input
              v-model.number="form.bigBlind"
              label="大盲注"
              type="number"
              :rules="blindRules"
              outlined
              class="blind-input"
              readonly
              :hint="`自动设置为小盲注的2倍`"
            >
              <template v-slot:prepend>
                <q-icon name="visibility" />
              </template>
            </q-input>
          </div>
        </div>

        <!-- 最大玩家数 -->
        <q-slider
          v-model="form.maxPlayers"
          :min="2"
          :max="6"
          :step="1"
          label
          :label-value="`最多 ${form.maxPlayers} 位玩家`"
          color="primary"
          class="form-field"
        />

        <!-- 房间类型 -->
        <q-toggle
          v-model="form.isPrivate"
          label="私人房间"
          color="primary"
          class="form-field"
        />

        <!-- 私人房间密码 -->
        <q-input
          v-if="form.isPrivate"
          v-model="form.password"
          label="房间密码"
          type="password"
          :rules="passwordRules"
          outlined
          class="form-field"
          hint="其他玩家需要密码才能加入"
        >
          <template v-slot:prepend>
            <q-icon name="lock" />
          </template>
        </q-input>

        <!-- 预览卡片 -->
        <q-card flat bordered class="preview-card">
          <q-card-section class="preview-header">
            <div class="text-subtitle2">房间预览</div>
          </q-card-section>
          
          <q-card-section class="preview-content">
            <div class="preview-grid">
              <div class="preview-item">
                <q-icon name="title" color="primary" />
                <span>{{ form.name || '未设置房间名' }}</span>
              </div>
              
              <div class="preview-item">
                <q-icon name="star" :color="chipLevelColor" />
                <span>{{ chipLevelText }}</span>
              </div>
              
              <div class="preview-item">
                <q-icon name="people" color="primary" />
                <span>最多 {{ form.maxPlayers }} 人</span>
              </div>
              
              <div class="preview-item">
                <q-icon name="account_balance" color="orange" />
                <span>最低 {{ form.minChips }} 筹码</span>
              </div>
              
              <div class="preview-item">
                <q-icon name="visibility" color="green" />
                <span>{{ form.smallBlind }}/{{ form.bigBlind }} 盲注</span>
              </div>
              
              <div class="preview-item">
                <q-icon :name="form.isPrivate ? 'lock' : 'public'" color="grey" />
                <span>{{ form.isPrivate ? '私人房间' : '公开房间' }}</span>
              </div>
            </div>
          </q-card-section>
        </q-card>
      </q-form>
    </q-card-section>

    <q-card-actions align="right" class="q-pt-none">
      <q-btn
        flat
        label="取消"
        color="grey"
        @click="$emit('close')"
      />
      
      <q-btn
        label="创建房间"
        color="primary"
        :loading="isLoading"
        @click="handleSubmit"
        unelevated
      >
        <template v-slot:loading>
          <q-spinner-hourglass class="on-left" />
          创建中...
        </template>
      </q-btn>
    </q-card-actions>
  </q-card>
</template>

<script>
import { defineComponent, ref, computed, reactive } from 'vue'
import { useQuasar } from 'quasar'
import { useGameStore } from '../stores/game'
import { useAuthStore } from '../stores/auth'

export default defineComponent({
  name: 'CreateRoomDialog',
  
  emits: ['close', 'created'],
  
  setup(props, { emit }) {
    const $q = useQuasar()
    const gameStore = useGameStore()
    const authStore = useAuthStore()
    
    const isLoading = ref(false)
    
    // 表单数据
    const form = reactive({
      name: '',
      chipLevel: 'low',
      minChips: 100,
      smallBlind: 5,
      bigBlind: 10,
      maxPlayers: 6,
      isPrivate: false,
      password: ''
    })
    
    // 筹码级别选项
    const chipLevelOptions = [
      { label: '低级场 (适合新手)', value: 'low' },
      { label: '中级场 (普通玩家)', value: 'medium' },
      { label: '高级场 (高手对决)', value: 'high' }
    ]
    
    // 计算属性
    const chipLevelColor = computed(() => {
      const colors = { low: 'green', medium: 'orange', high: 'red' }
      return colors[form.chipLevel]
    })
    
    const chipLevelText = computed(() => {
      const option = chipLevelOptions.find(opt => opt.value === form.chipLevel)
      return option ? option.label : '未知级别'
    })
    
    const recommendedMinChips = computed(() => {
      return form.bigBlind * 20 // 建议最低筹码为大盲注的20倍
    })
    
    // 验证规则
    const nameRules = [
      val => !!val || '请输入房间名称',
      val => val.length >= 2 || '房间名称至少2个字符',
      val => val.length <= 50 || '房间名称不能超过50个字符'
    ]
    
    const minChipsRules = [
      val => !!val || '请设置最低筹码要求',
      val => val >= form.bigBlind * 10 || `最低筹码不能少于大盲注的10倍 (${form.bigBlind * 10})`,
      val => val <= authStore.userChips || `最低筹码不能超过您的筹码数 (${authStore.userChips})`
    ]
    
    const blindRules = [
      val => !!val || '请设置盲注',
      val => val > 0 || '盲注必须大于0'
    ]
    
    const passwordRules = [
      val => !form.isPrivate || !!val || '私人房间需要设置密码',
      val => !form.isPrivate || val.length >= 4 || '密码至少4个字符'
    ]
    
    // 方法
    const updateChipSettings = (chipLevel) => {
      // 根据筹码级别自动调整默认值
      const settings = {
        low: { minChips: 100, smallBlind: 5 },
        medium: { minChips: 500, smallBlind: 25 },
        high: { minChips: 2000, smallBlind: 100 }
      }
      
      const setting = settings[chipLevel]
      if (setting) {
        form.minChips = setting.minChips
        form.smallBlind = setting.smallBlind
        form.bigBlind = setting.smallBlind * 2
      }
    }
    
    const updateBigBlind = () => {
      // 大盲注自动设为小盲注的2倍
      form.bigBlind = form.smallBlind * 2
    }
    
    const validateForm = () => {
      // 检查所有验证规则
      const nameValid = nameRules.every(rule => rule(form.name) === true)
      const minChipsValid = minChipsRules.every(rule => rule(form.minChips) === true)
      const blindValid = blindRules.every(rule => rule(form.smallBlind) === true)
      const passwordValid = passwordRules.every(rule => rule(form.password) === true)
      
      return nameValid && minChipsValid && blindValid && passwordValid
    }
    
    const handleSubmit = async () => {
      if (!validateForm()) {
        $q.notify({
          type: 'negative',
          message: '请检查表单输入',
          position: 'top'
        })
        return
      }
      
      // 检查用户筹码是否足够
      if (authStore.userChips < form.minChips) {
        $q.notify({
          type: 'negative',
          message: '您的筹码不足以创建这个级别的房间',
          position: 'top'
        })
        return
      }
      
      isLoading.value = true
      
      try {
        const roomData = {
          name: form.name,
          chip_level: form.chipLevel,
          min_chips: form.minChips,
          small_blind: form.smallBlind,
          big_blind: form.bigBlind,
          max_players: form.maxPlayers,
          is_private: form.isPrivate,
          password: form.isPrivate ? form.password : null
        }
        
        const result = await gameStore.createRoom(roomData)
        
        if (result.success) {
          $q.notify({
            type: 'positive',
            message: '房间创建成功！',
            position: 'top'
          })
          
          emit('created', result.room)
        } else {
          $q.notify({
            type: 'negative',
            message: result.message,
            position: 'top'
          })
        }
      } catch (error) {
        $q.notify({
          type: 'negative',
          message: '创建房间失败，请稍后重试',
          position: 'top'
        })
      } finally {
        isLoading.value = false
      }
    }
    
    return {
      form,
      isLoading,
      chipLevelOptions,
      chipLevelColor,
      chipLevelText,
      recommendedMinChips,
      nameRules,
      minChipsRules,
      blindRules,
      passwordRules,
      updateChipSettings,
      updateBigBlind,
      handleSubmit
    }
  }
})
</script>

<style scoped lang="scss">
.create-room-form {
  .form-field {
    margin-bottom: 1rem;
  }
  
  .blind-settings {
    margin-bottom: 1rem;
    
    .blind-title {
      font-size: 0.875rem;
      font-weight: 500;
      color: #666;
      margin-bottom: 0.5rem;
    }
    
    .blind-inputs {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 1rem;
    }
  }
}

.preview-card {
  margin-top: 1rem;
  background: #f8f9fa;
  
  .preview-header {
    padding: 0.75rem 1rem 0.25rem 1rem;
    
    .text-subtitle2 {
      color: #1976d2;
      font-weight: 500;
    }
  }
  
  .preview-content {
    padding: 0.5rem 1rem 1rem 1rem;
    
    .preview-grid {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 0.75rem;
    }
    
    .preview-item {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      font-size: 0.875rem;
      color: #333;
    }
  }
}

// 响应式设计
@media (max-width: 480px) {
  .blind-inputs {
    grid-template-columns: 1fr !important;
  }
  
  .preview-grid {
    grid-template-columns: 1fr !important;
  }
}
</style> 