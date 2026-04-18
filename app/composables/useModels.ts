export const useModels = () => {
    const supabase = useSupabaseClient();

    const fetchModel = async (id: number) => {
        return supabase
            .from("models")
            .select("*, providers(id, name, npm, env, doc, api)")
            .eq("id", id)
            .single();
    };

    return { fetchModel };
};
