import {Button,Card,Checkbox,Col,Form,Input,Row,Select,Slider,Tag,Typography,message} from 'antd';
import {useEffect,useState} from 'react';
import {deleteJournal,listJournalPrompts,listJournals,saveJournal,updateJournal} from '../api/journal';
import {EmptyState} from '../components/common/EmptyState';
import {MoodCard} from '../components/common/MoodCard';
import {JOURNAL_PROMPT_FALLBACK,JOURNAL_PROMPT_SELECTED,JOURNAL_PROMPT_TITLE} from '../constants/journal';
import type {Journal,JournalPrompts} from '../types';

type JournalFormValues={title:string;content:string;mood_level:number;weather:string;prompt:string;is_private:boolean};
const DEFAULTS:JournalFormValues={title:'',content:'',mood_level:7,weather:'晴',prompt:'',is_private:true};

export function Journals(){
  const [items,setItems]=useState<Journal[]>([]);
  const [form]=Form.useForm<JournalFormValues>();
  const level=Form.useWatch('mood_level',form)??DEFAULTS.mood_level;
  const selectedPrompt=Form.useWatch('prompt',form)??'';
  const [promptData,setPromptData]=useState<JournalPrompts|null>(null);
  const [promptsFailed,setPromptsFailed]=useState(false);
  const [editing,setEditing]=useState<Journal|null>(null);
  const load=()=>listJournals().then(setItems).catch(e=>message.error(e.message));
  useEffect(()=>{load()},[]);
  // 心情指数变化时拉取对应的三条书写提示；加载失败只隐藏提示区，正文始终可自由书写
  useEffect(()=>{let stale=false;setPromptsFailed(false);
    listJournalPrompts(level).then(d=>{if(!stale)setPromptData(d)}).catch(()=>{if(!stale){setPromptData(null);setPromptsFailed(true)}});
    return()=>{stale=true}},[level]);
  // 点一条提示：问题带进正文（可再改写），并记住这篇日记当时选的是哪条
  const applyPrompt=(q:string)=>{
    const prev=form.getFieldValue('prompt')||'';
    const content=form.getFieldValue('content')||'';
    let next:string;
    if(prev&&content.includes(prev))next=content.replace(prev,q);
    else if(!content.trim())next=q+'\n';
    else next=content.replace(/\s+$/,'')+'\n'+q+'\n';
    form.setFieldsValue({content:next,prompt:q});
  };
  const submit=async(v:JournalFormValues)=>{
    try{
      if(editing){await updateJournal(editing.id,v);message.success('日记已更新')}
      else{await saveJournal(v);message.success('日记已保存')}
      setEditing(null);form.resetFields();load();
    }catch(e){message.error((e as Error).message)}
  };
  const startEdit=(j:Journal)=>{setEditing(j);form.setFieldsValue({title:j.title,content:j.content,mood_level:j.mood_level,weather:j.weather,prompt:j.prompt||'',is_private:j.is_private})};
  const remove=async(j:Journal)=>{
    try{
      await deleteJournal(j.id);message.success('已删除');
      if(editing?.id===j.id){setEditing(null);form.resetFields()}
      load();
    }catch(e){message.error((e as Error).message)}
  };
  return <>
    <Typography.Title>日记本</Typography.Title>
    <Row gutter={[20,20]}>
      <Col xs={24} lg={10}>
        <Card title={editing?`编辑日记 #${editing.id}`:'写一页日记'}>
          <Form form={form} layout="vertical" initialValues={DEFAULTS} onFinish={submit}>
            <Form.Item name="title" label="标题" rules={[{required:true}]}><Input placeholder="给今天一句标题"/></Form.Item>
            <Form.Item name="mood_level" label="心情指数"><Slider min={1} max={10}/></Form.Item>
            {promptData&&promptData.prompts.length>0&&<div className="prompt-block">
              <div className="muted prompt-title">{JOURNAL_PROMPT_TITLE}{promptData.band?`（${promptData.band}）`:''}</div>
              <div className="prompt-list">{promptData.prompts.map(q=><button type="button" key={q} className={`prompt-item${q===selectedPrompt?' prompt-item--active':''}`} onClick={()=>applyPrompt(q)}>{q}</button>)}</div>
            </div>}
            {promptsFailed&&<div className="prompt-fallback">{JOURNAL_PROMPT_FALLBACK}</div>}
            <Form.Item name="content" label="正文（支持 Markdown 风格文本）" rules={[{required:true}]}><Input.TextArea rows={8} placeholder="今天发生了什么？"/></Form.Item>
            <Form.Item name="prompt" hidden><Input/></Form.Item>
            <Form.Item name="weather" label="天气"><Select options={['晴','阴','雨','风'].map(v=>({value:v}))}/></Form.Item>
            <Form.Item name="is_private" valuePropName="checked"><Checkbox>仅自己可见</Checkbox></Form.Item>
            <Button type="primary" htmlType="submit">{editing?'保存修改':'保存日记'}</Button>
            {editing&&<Button style={{marginLeft:8}} onClick={()=>{setEditing(null);form.resetFields()}}>取消编辑</Button>}
          </Form>
        </Card>
      </Col>
      <Col xs={24} lg={14}>
        <Card title="时间轴">{items.length?<div className="timeline">{items.map(j=><article className="journal" key={j.id}>
          <div><Tag color="purple">{j.created_at.slice(0,10)}</Tag><b>{j.title}</b><span className="muted"> · {j.weather} · 心情 {j.mood_level}/10</span></div>
          {j.prompt&&<div className="journal-prompt"><Tag color="green">{JOURNAL_PROMPT_SELECTED}</Tag><span>{j.prompt}</span></div>}
          <p>{j.content}</p>
          <Button size="small" onClick={()=>startEdit(j)}>编辑</Button>
          <Button danger size="small" style={{marginLeft:8}} onClick={()=>remove(j)}>删除</Button>
          <div className="journal-mood"><MoodCard mood={{id:j.id,user_id:0,mood_level:j.mood_level,mood_tags:'["calm"]',note:'日记中的情绪',record_date:j.created_at,created_at:j.created_at}}/></div>
        </article>)}</div>:<EmptyState title="写下第一篇日记，和自己好好聊聊"/>}</Card>
      </Col>
    </Row>
  </>;
}
