<script lang="ts">
  import StatusBanner from "$shared/components/StatusBanner.svelte";
  import PageTitle from "$shared/components/PageTitle.svelte";

  import { t } from "$shared/stores/locale.svelte";
  import ConfirmDialog from "$shared/components/ConfirmDialog.svelte";
  import { Button } from "$shared/components/buttons";
  import ProfileCard from "./components/ProfileCard.svelte";
  import ProfileForm from "./components/ProfileForm.svelte";
  import { createProfilesStore } from "./store.svelte";

  let { data } = $props();
  const store = createProfilesStore();
</script>

<div class="space-y-6">
  <!-- Top Header -->
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
    <PageTitle title={t("profiles.title")} />

    <Button variant="primary" size="sm" onclick={store.openCreate}>
      {t("profiles.new_profile")}
    </Button>
  </div>

  <!-- Status Alerts -->
  <StatusBanner />

  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
    {#each data.profiles as profile (profile.id)}
      <ProfileCard
        {profile}
        active={data.activeProfile?.id === profile.id}
        canDelete={data.profiles.length > 1}
        onSelect={() => store.select(profile.id)}
        onEdit={() => store.openEdit(profile)}
        onDelete={() => store.requestDelete(profile.id, profile.name)}
      />
    {/each}
  </div>
</div>

<!-- Modal: Add/Edit Profile -->
<ProfileForm
  bind:show={store.showModal}
  id={store.form.id}
  bind:name={store.form.name}
  onSave={store.save}
/>

<!-- Profile Delete Confirmation Modal -->
<ConfirmDialog
  bind:show={store.showDeleteConfirm}
  title="Remover Perfil"
  message={`Tem certeza de que deseja remover o perfil '${store.profileToDelete?.name || ""}'?\nEssa ação não poderá ser desfeita.`}
  confirmText="Remover Perfil"
  type="danger"
  onConfirm={store.confirmDelete}
/>
