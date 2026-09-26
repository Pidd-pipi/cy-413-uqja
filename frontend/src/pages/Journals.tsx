import {Button,Card,Checkbox,Col,Form,Input,Row,Select,Slider,Tag,Typography,message} from 'antd';
import {useEffect,useState} from 'react';
import {deleteJournal,listJournalPrompts,listJournals,saveJournal} from '../api/journal';
import {EmptyState} from '../components/common/EmptyState';
import {MoodCard} from '../components/common/MoodCard';
import type {Journal} from '../types';

type JournalFormValues={title:string;content:string;mood_level:number;weather:string;is_private:boolean};

export function Journals(){
  const [items,setItems]=useState<Journal[]>([]);
  const [form]=Form.useForm<JournalFormValues>();
  const moodLevel=Form.useWatch('mood_level',form)??7;
  const [prompts,setPrompts]=useState<string[]>([]);
  const [promptsFailed,setPromptsFailed]=useState(false);
  const [chosenPrompt,setChosenPrompt]=useState('');
  const load=()=>listJournals().then(setItems).catch(e=>message.error(e.message));
  useEffect(()=>{load()},[]);
  useEffect(()=>{
    let alive=true;
    const timer=setTimeout(()=>{
      listJournalPrompts(moodLevel)
        .then(p=>{if(alive){setPrompts(p);setPromptsFailed(false)}})
        .catch(()=>{if(alive){setPrompts([]);setPromptsFailed(true)}});
    },250);
    return ()=>{alive=false;clearTimeout(timer)};
  },[moodLevel]);
  const usePrompt=(p:string)=>{
    const cur=form.getFieldValue('content')||'';
    if(!cur.includes(p))form.setFieldValue('content',cur?`${cur}\n${p}`:`${p}\n`);
    setChosenPrompt(p);
  };
  const submit=async(v:JournalFormValues)=>{
    try{
      await saveJournal({...v,prompt:chosenPrompt});
      message.success('日记已保存');
      setChosenPrompt('');
      form.resetFields();
      load();
    }catch(e){message.error((e as Error).message)}
  };
  return <><Typography.Title>日记本</Typography.Title><Row gutter={[20,20]}><Col xs={24} lg={10}><Card title="写一页日记"><Form form={form} layout="vertical" initialValues={{mood_level:7,is_private:true,weather:'晴'}} onFinish={submit}><Form.Item name="title" label="标题" rules={[{required:true}]}><Input placeholder="给今天一句标题"/></Form.Item><Form.Item name="mood_level" label="心情指数"><Slider min={1} max={10}/></Form.Item><div className="prompt-block"><span className="prompt-block__title">不知道从哪写起？点一条带进正文，可以自由改写</span>{promptsFailed?<span className="muted">书写提示加载失败，没关系，直接自由书写就好。</span>:prompts.length?<div className="prompt-list">{prompts.map(p=><button type="button" key={p} className={p===chosenPrompt?'prompt-item prompt-item--active':'prompt-item'} onClick={()=>usePrompt(p)}>{p}</button>)}</div>:<span className="muted">正在为你准备书写提示…</span>}</div><Form.Item name="content" label="正文（支持 Markdown 风格文本）" rules={[{required:true}]}><Input.TextArea rows={8} placeholder="今天发生了什么？"/></Form.Item><Form.Item name="weather" label="天气"><Select options={['晴','阴','雨','风'].map(v=>({value:v}))}/></Form.Item><Form.Item name="is_private" valuePropName="checked"><Checkbox>仅自己可见</Checkbox></Form.Item><Button type="primary" htmlType="submit">保存日记</Button></Form></Card></Col><Col xs={24} lg={14}><Card title="时间轴">{items.length?<div className="timeline">{items.map(j=><article className="journal" key={j.id}><div><Tag color="purple">{j.created_at.slice(0,10)}</Tag><b>{j.title}</b><span className="muted"> · {j.weather} · 心情 {j.mood_level}/10</span></div>{j.prompt&&<div className="journal-prompt"><b>书写提示</b>{j.prompt}</div>}<p>{j.content}</p><Button danger size="small" onClick={async()=>{await deleteJournal(j.id);message.success('已删除');load()}}>删除</Button><div className="journal-mood"><MoodCard mood={{id:j.id,user_id:0,mood_level:j.mood_level,mood_tags:'["calm"]',note:'日记中的情绪',record_date:j.created_at,created_at:j.created_at}}/></div></article>)}</div>:<EmptyState title="写下第一篇日记，和自己好好聊聊"/>}</Card></Col></Row></>;
}
